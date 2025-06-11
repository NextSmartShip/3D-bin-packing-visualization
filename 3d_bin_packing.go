package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

type Item struct {
	Length int32
	Width  int32
	Height int32
	Qty    int32
}

type Container struct {
	Length int32
	Width  int32
	Height int32
}

// Position表示物体在3D空间中的位置
type Position struct {
	X int32
	Y int32
	Z int32
}

// PlacedItem表示已放置的物品
type PlacedItem struct {
	Length   int32
	Width    int32
	Height   int32
	Position Position
}

// ItemOrientation表示物品的一种摆放方向
type ItemOrientation struct {
	Length int32
	Width  int32
	Height int32
}

// ItemPlacementRecord 记录物品的摆放信息
type ItemPlacementRecord struct {
	ItemIndex   int             // 物品在Items数组中的索引
	Position    Position        // 摆放位置
	Orientation ItemOrientation // 摆放方向
}

// PackingSession 封装一次装箱会话中的共享状态
type PackingSession struct {
	Container            Container
	Items                []*Item
	StartTime            time.Time
	WarningExecutionTime time.Duration
	UsedPositions        map[Position]bool
	// ItemPlacements 记录每个物品的摆放信息
	ItemPlacements []ItemPlacementRecord
}

// PlacementState 封装当前放置的状态
type PlacementState struct {
	Index          int
	PlacedItems    []PlacedItem
	Orientation    ItemOrientation
	Position       Position
	OrientIdx      int
	TotalOrients   int
	PosIdx         int
	TotalPositions int
}

// RecursivePackingParams 封装递归装箱函数的参数
type RecursivePackingParams struct {
	Session       *PackingSession
	State         *PlacementState
	IsNearTimeout bool
}

// BinPacker 装箱器，封装装箱逻辑
type BinPacker struct {
	session *PackingSession
	options *PackOptions
}

// 默认警戒时间，超过此时间后将减少回溯
const defaultWarningExecutionTime = 3 * time.Second

// PackOptions 打包选项
type PackOptions struct {
	// 警戒时间，超过此时间后将减少回溯
	WarningExecutionTime time.Duration
	// 物品排序函数，用于决定物品的放置顺序
	ItemSortFunc func([]*Item)
}

// PackOption 定义打包选项的设置函数类型
type PackOption func(*PackOptions)

// WithWarningExecutionTime 设置警戒时间选项
func WithWarningExecutionTime(t time.Duration) PackOption {
	return func(o *PackOptions) {
		o.WarningExecutionTime = t
	}
}

// WithItemSortFunc 设置物品排序方法
func WithItemSortFunc(sortFunc func([]*Item)) PackOption {
	return func(o *PackOptions) {
		o.ItemSortFunc = sortFunc
	}
}

// WithBaseAreaSortFunc 定义按底面积排序的选项，底面积大的物品优先放置
// 物品将按底面积从大到小排序
// 底面积相同的物品将按体积从大到小排序
func WithBaseAreaSortFunc() PackOption {
	return WithItemSortFunc(func(items []*Item) {
		sort.Slice(items, func(i, j int) bool {
			// 计算物品i的最大底面积
			itemI := items[i]
			baseAreaI1 := itemI.Length * itemI.Width  // 长边作为底面
			baseAreaI2 := itemI.Length * itemI.Height // 高边作为底面
			baseAreaI3 := itemI.Width * itemI.Height  // 宽边作为底面
			maxBaseAreaI := maxInt32(baseAreaI1, baseAreaI2, baseAreaI3)

			// 计算物品j的最大底面积
			itemJ := items[j]
			baseAreaJ1 := itemJ.Length * itemJ.Width  // 长边作为底面
			baseAreaJ2 := itemJ.Length * itemJ.Height // 高边作为底面
			baseAreaJ3 := itemJ.Width * itemJ.Height  // 宽边作为底面
			maxBaseAreaJ := maxInt32(baseAreaJ1, baseAreaJ2, baseAreaJ3)

			// 如果底面积相同，则按体积排序
			if maxBaseAreaI == maxBaseAreaJ {
				volI := int64(itemI.Length) * int64(itemI.Width) * int64(itemI.Height)
				volJ := int64(itemJ.Length) * int64(itemJ.Width) * int64(itemJ.Height)
				return volI > volJ // 体积大的优先放置
			}

			return maxBaseAreaI > maxBaseAreaJ // 底面积大的优先放置
		})
	})
}

// WithDimensionSortFunc 返回按物品尺寸维度排序的选项
// 首先按最长边从大到小排序，相同时按次长边从大到小排序，再相同时按最短边排序
// 这种排序方式适合长条形物品较多的情况，有助于更有效地利用空间
func WithDimensionSortFunc() PackOption {
	return WithItemSortFunc(func(items []*Item) {
		sort.Slice(items, func(i, j int) bool {
			// 获取物品i的三个维度并按从大到小排序
			dimsI := []int32{items[i].Length, items[i].Width, items[i].Height}
			sort.Slice(dimsI, func(a, b int) bool {
				return dimsI[a] > dimsI[b]
			})

			// 获取物品j的三个维度并按从大到小排序
			dimsJ := []int32{items[j].Length, items[j].Width, items[j].Height}
			sort.Slice(dimsJ, func(a, b int) bool {
				return dimsJ[a] > dimsJ[b]
			})

			// 依次比较最长边、次长边、最短边
			for k := 0; k < 3; k++ {
				if dimsI[k] != dimsJ[k] {
					return dimsI[k] > dimsJ[k] // 维度大的排在前面
				}
			}

			// 如果所有维度都相等，则返回false（保持原顺序）
			return false
		})
	})
}

// maxInt32 返回多个int32中的最大值
func maxInt32(values ...int32) int32 {
	if len(values) == 0 {
		return 0
	}

	maxVal := values[0]
	for _, v := range values[1:] {
		if v > maxVal {
			maxVal = v
		}
	}

	return maxVal
}

// 默认的物品排序函数 - 按体积从大到小排序
func defaultItemSortFunc(items []*Item) {
	sort.Slice(items, func(i, j int) bool {
		volI := items[i].Length * items[i].Width * items[i].Height
		volJ := items[j].Length * items[j].Width * items[j].Height
		return volI > volJ
	})
}

// CanPack检查物品是否能装入容器，使用可变参数传入选项
func CanPack(container Container, items []*Item, opts ...PackOption) bool {
	// 创建默认选项
	options := &PackOptions{
		WarningExecutionTime: defaultWarningExecutionTime,
		ItemSortFunc:         defaultItemSortFunc,
	}

	// 应用所有选项
	for _, opt := range opts {
		opt(options)
	}

	// 使用新的BinPacker实现
	packer := NewBinPacker(container, options)
	success := packer.Pack(items)

	if success {
		// 生成可视图
		generatePackingVisualization(packer, container)
	}

	return success
}

// NewBinPacker 创建新的装箱器
func NewBinPacker(container Container, options *PackOptions) *BinPacker {
	if options == nil {
		options = &PackOptions{
			WarningExecutionTime: defaultWarningExecutionTime,
			ItemSortFunc:         defaultItemSortFunc,
		}
	}

	return &BinPacker{
		session: &PackingSession{
			Container:            container,
			WarningExecutionTime: options.WarningExecutionTime,
			UsedPositions:        make(map[Position]bool),
			ItemPlacements:       make([]ItemPlacementRecord, 0),
		},
		options: options,
	}
}

// Pack 执行装箱，判断是否能装下所有物品
func (p *BinPacker) Pack(items []*Item) bool {
	// 首先计算所有物品的总体积
	var totalVolume int64
	for _, item := range items {
		// 计算单个物品的体积并乘以数量
		itemVolume := int64(item.Length) * int64(item.Width) * int64(item.Height) * int64(item.Qty)
		totalVolume += itemVolume
	}

	// 计算容器的体积
	containerVolume := int64(p.session.Container.Length) * int64(p.session.Container.Width) * int64(p.session.Container.Height)

	// 如果物品总体积大于容器体积，肯定无法装入
	if totalVolume > containerVolume {
		return false
	}

	// 若物品和容器都是空的，返回true
	if len(items) == 0 || (p.session.Container.Length == 0 && p.session.Container.Width == 0 && p.session.Container.Height == 0) {
		return true
	}

	// 对物品进行排序，启发式方法，先放置大物品
	// 创建完整的物品列表，考虑每个物品的数量
	var fullItemList []*Item
	for _, item := range items {
		for i := int32(0); i < item.Qty; i++ {
			// 复制一份物品，数量设为1
			newItem := &Item{
				Length: item.Length,
				Width:  item.Width,
				Height: item.Height,
				Qty:    1,
			}
			fullItemList = append(fullItemList, newItem)
		}
	}

	// 使用指定的排序函数对物品排序
	p.options.ItemSortFunc(fullItemList)

	// 更新session中的items和startTime
	p.session.Items = fullItemList
	p.session.StartTime = time.Now()

	// 重置已使用位置和摆放记录
	p.session.UsedPositions = make(map[Position]bool)
	p.session.ItemPlacements = make([]ItemPlacementRecord, 0)

	// 开始递归装箱
	return p.doPackItems(0, []PlacedItem{})
}

// GetItemPlacements 获取物品摆放记录
func (p *BinPacker) GetItemPlacements() []ItemPlacementRecord {
	return p.session.ItemPlacements
}

// doPackItems BinPacker的核心递归装箱方法
func (p *BinPacker) doPackItems(index int, placedItems []PlacedItem) bool {
	// 所有物品都已放置，成功
	if index >= len(p.session.Items) {
		return true
	}

	// 获取当前物品和其可能的方向
	currentItem := p.session.Items[index]
	orientations := generateOrientations(currentItem)

	// 统一使用WarningExecutionTime来判断是否接近超时
	elapsed := time.Since(p.session.StartTime)
	isNearTimeout := elapsed > p.session.WarningExecutionTime/2 // 达到一半警戒时间时开始减少尝试

	// 如果接近超时，减少尝试的方向数量
	if isNearTimeout && len(orientations) > 2 {
		orientations = orientations[:2] // 只尝试前2个方向
	}

	// 尝试每种旋转方向
	for orientIdx, orientation := range orientations {
		if !canFitInContainer(orientation, p.session.Container) {
			continue
		}

		state := &PlacementState{
			Index:        index,
			PlacedItems:  placedItems,
			Orientation:  orientation,
			OrientIdx:    orientIdx,
			TotalOrients: len(orientations),
		}

		params := &RecursivePackingParams{
			Session:       p.session,
			State:         state,
			IsNearTimeout: isNearTimeout,
		}

		if p.tryPlaceItemInOrientation(params) {
			return true
		}

		// 如果已接近超时，不再尝试其他方向
		if isNearTimeout {
			break
		}
	}

	return false
}

// tryPlaceItemInOrientation 尝试在特定方向放置物品
func (p *BinPacker) tryPlaceItemInOrientation(params *RecursivePackingParams) bool {
	// 获取所有可能的放置点
	possiblePositions := p.findPossiblePositions(params.State.Orientation, params.State.PlacedItems)

	// 接近超时时，只尝试有限的几个位置
	if params.IsNearTimeout && len(possiblePositions) > 3 {
		possiblePositions = possiblePositions[:3]
	}

	// 更新状态中的总位置数
	params.State.TotalPositions = len(possiblePositions)

	// 尝试每个可能的位置
	for posIdx, pos := range possiblePositions {
		params.State.Position = pos
		params.State.PosIdx = posIdx

		if p.tryPlaceItemAtPosition(params) {
			return true
		}
	}

	return false
}

// tryPlaceItemAtPosition 尝试在特定位置放置物品
func (p *BinPacker) tryPlaceItemAtPosition(params *RecursivePackingParams) bool {
	pos := params.State.Position

	// 标记当前位置为已使用
	params.Session.UsedPositions[pos] = true
	defer func() {
		// 回溯时，取消标记当前位置为已使用
		params.Session.UsedPositions[pos] = false
	}()

	// 创建新的已放置物品
	newPlacedItem := PlacedItem{
		Length:   params.State.Orientation.Length,
		Width:    params.State.Orientation.Width,
		Height:   params.State.Orientation.Height,
		Position: pos,
	}

	// 检查是否与已放置物品冲突
	if !hasCollision(newPlacedItem, params.State.PlacedItems) {
		// 放置物品 - 创建新的slice而不是修改原slice，这是有意的
		newPlacedItems := append(params.State.PlacedItems, newPlacedItem) //nolint:gocritic // intentionally creating new slice

		// 递归放置下一个物品
		if p.doPackItems(params.State.Index+1, newPlacedItems) {
			// 成功放置物品，记录摆放信息
			placementRecord := ItemPlacementRecord{
				ItemIndex:   params.State.Index,
				Position:    pos,
				Orientation: params.State.Orientation,
			}
			params.Session.ItemPlacements = append(params.Session.ItemPlacements, placementRecord)
			return true
		}

		// 检查是否应该提前退出
		if p.shouldExitEarly(params) {
			return false
		}
	}

	return false
}

// shouldExitEarly 检查是否应该提前退出
func (p *BinPacker) shouldExitEarly(params *RecursivePackingParams) bool {
	elapsed := time.Since(params.Session.StartTime)

	// 如果超过警戒时间，更积极地退出
	if elapsed > params.Session.WarningExecutionTime {
		// 如果已经尝试了几个位置或方向，就退出
		if params.State.PosIdx >= 1 || params.State.OrientIdx >= 1 {
			return true
		}
	}

	// 如果超过警戒时间很多，直接退出
	if elapsed > params.Session.WarningExecutionTime*2 {
		return true
	}

	return false
}

// findPossiblePositions 找出所有可能的放置位置
func (p *BinPacker) findPossiblePositions(itemSize ItemOrientation, placedItems []PlacedItem) []Position {
	// 如果没有已放置的物品，只有容器原点可用
	if len(placedItems) == 0 {
		// 如果原点正在被使用，返回空列表
		if p.session.UsedPositions[Position{0, 0, 0}] {
			return []Position{}
		}
		return []Position{{0, 0, 0}}
	}

	// 候选位置集合，无需使用map，直接使用切片
	var candidatePositions []Position

	// 添加容器原点，如果没有被使用
	if !p.session.UsedPositions[Position{0, 0, 0}] {
		candidatePositions = append(candidatePositions, Position{0, 0, 0})
	}

	// 添加每个已放置物品的"极限点"
	for _, placed := range placedItems {
		// 添加物品顶部的位置（适合物品堆叠）
		topPosition := Position{
			placed.Position.X,
			placed.Position.Y,
			placed.Position.Z + placed.Height,
		}
		// 只添加未被使用的位置
		if !p.session.UsedPositions[topPosition] {
			candidatePositions = append(candidatePositions, topPosition)
		}

		// 添加物品右侧的位置
		rightPosition := Position{
			placed.Position.X + placed.Length,
			placed.Position.Y,
			placed.Position.Z,
		}
		// 只添加未被使用的位置
		if !p.session.UsedPositions[rightPosition] {
			candidatePositions = append(candidatePositions, rightPosition)
		}

		// 添加物品前方的位置
		frontPosition := Position{
			placed.Position.X,
			placed.Position.Y + placed.Width,
			placed.Position.Z,
		}
		// 只添加未被使用的位置
		if !p.session.UsedPositions[frontPosition] {
			candidatePositions = append(candidatePositions, frontPosition)
		}
	}

	// 过滤出有效的位置
	var validPositions []Position
	for _, pos := range candidatePositions {
		// 检查位置是否在容器范围内
		if pos.X >= 0 && pos.Y >= 0 && pos.Z >= 0 &&
			pos.X+itemSize.Length <= p.session.Container.Length &&
			pos.Y+itemSize.Width <= p.session.Container.Width &&
			pos.Z+itemSize.Height <= p.session.Container.Height {
			validPositions = append(validPositions, pos)
		}
	}

	return validPositions
}

// 检查物品是否与已放置物品有碰撞
func hasCollision(item PlacedItem, placedItems []PlacedItem) bool {
	for _, placed := range placedItems {
		// 检查两个3D矩形是否重叠
		if item.Position.X+item.Length > placed.Position.X &&
			placed.Position.X+placed.Length > item.Position.X &&
			item.Position.Y+item.Width > placed.Position.Y &&
			placed.Position.Y+placed.Width > item.Position.Y &&
			item.Position.Z+item.Height > placed.Position.Z &&
			placed.Position.Z+placed.Height > item.Position.Z {
			return true // 有碰撞
		}
	}
	return false // 无碰撞
}

// 检查物品在某个方向是否能放入容器
func canFitInContainer(orientation ItemOrientation, container Container) bool {
	return orientation.Length <= container.Length &&
		orientation.Width <= container.Width &&
		orientation.Height <= container.Height
}

// 生成物品的所有可能旋转方向
func generateOrientations(item *Item) []ItemOrientation {
	length := item.Length
	width := item.Width
	height := item.Height

	// 如果长方体的某两边相等，则旋转方向可能相同。使用Set来避免重复的方向
	orientationSet := make(map[ItemOrientation]bool)

	// 生成所有6种可能的旋转方向
	orientationSet[ItemOrientation{length, width, height}] = true // 原始方向
	orientationSet[ItemOrientation{length, height, width}] = true // 绕X轴旋转90度
	orientationSet[ItemOrientation{width, length, height}] = true // 绕Z轴旋转90度
	orientationSet[ItemOrientation{width, height, length}] = true // 复合旋转1
	orientationSet[ItemOrientation{height, length, width}] = true // 复合旋转2
	orientationSet[ItemOrientation{height, width, length}] = true // 复合旋转3

	// 将Set转换为切片
	var orientations []ItemOrientation
	for orientation := range orientationSet {
		orientations = append(orientations, orientation)
	}

	return orientations
}

// generatePackingVisualization 生成装箱3D可视化数据
func generatePackingVisualization(packer *BinPacker, container Container) {
	placements := packer.GetItemPlacements()

	fmt.Printf("\n==================== 3D Bin Packing Visualization ====================\n")
	fmt.Printf("Container Dimensions: %d×%d×%d (L×W×H)\n", container.Length, container.Width, container.Height)
	fmt.Printf("Total Items Placed: %d\n", len(placements))
	fmt.Printf("Generating 3D visualization data for web viewer...\n")
	fmt.Printf("======================================================================\n\n")

	// 生成3D可视化JSON数据
	err := generate3DVisualizationJSON(placements, container, "bin_packing_3d.json")
	if err != nil {
		fmt.Printf("Error generating 3D visualization data: %v\n", err)
		return
	}

	fmt.Printf("3D visualization data saved as 'bin_packing_3d.json'\n")
	fmt.Printf("Open 'bin_packing_viewer.html' in your browser to view the interactive 3D visualization\n")

	// 生成空间利用率统计
	generateSpaceUtilizationStats(placements, container)

	// 生成HTML查看器
	generateHTMLViewer()
}

// Visualization3DData 3D可视化数据结构
type Visualization3DData struct {
	Container Container3D `json:"container"`
	Items     []Item3D    `json:"items"`
	Stats     Stats3D     `json:"stats"`
}

// Container3D 容器3D数据
type Container3D struct {
	Length int32 `json:"length"`
	Width  int32 `json:"width"`
	Height int32 `json:"height"`
}

// Item3D 物品3D数据
type Item3D struct {
	ID         int       `json:"id"`
	Position   Position  `json:"position"`
	Dimensions Dimension `json:"dimensions"`
	Color      string    `json:"color"`
}

// Dimension 尺寸信息
type Dimension struct {
	Length int32 `json:"length"`
	Width  int32 `json:"width"`
	Height int32 `json:"height"`
}

// Stats3D 统计信息
type Stats3D struct {
	TotalItems      int     `json:"totalItems"`
	ContainerVolume int64   `json:"containerVolume"`
	ItemsVolume     int64   `json:"itemsVolume"`
	UtilizationRate float64 `json:"utilizationRate"`
}

// generate3DVisualizationJSON 生成3D可视化JSON数据
func generate3DVisualizationJSON(placements []ItemPlacementRecord, container Container, filename string) error {
	// 预定义颜色数组
	colors := []string{
		"#ff6464", // 红色
		"#64ff64", // 绿色
		"#6464ff", // 蓝色
		"#ffff64", // 黄色
		"#ff64ff", // 品红
		"#64ffff", // 青色
		"#ff9664", // 橙色
		"#9664ff", // 紫色
		"#64ff96", // 浅绿
		"#ffc896", // 桃色
	}

	// 转换物品数据
	items := make([]Item3D, len(placements))
	var totalItemVolume int64

	for i, placement := range placements {
		items[i] = Item3D{
			ID:       i + 1,
			Position: placement.Position,
			Dimensions: Dimension{
				Length: placement.Orientation.Length,
				Width:  placement.Orientation.Width,
				Height: placement.Orientation.Height,
			},
			Color: colors[i%len(colors)],
		}

		// 计算体积
		itemVolume := int64(placement.Orientation.Length) *
			int64(placement.Orientation.Width) *
			int64(placement.Orientation.Height)
		totalItemVolume += itemVolume
	}

	// 计算统计信息
	containerVolume := int64(container.Length) * int64(container.Width) * int64(container.Height)
	utilizationRate := float64(totalItemVolume) / float64(containerVolume) * 100

	// 创建可视化数据
	vizData := Visualization3DData{
		Container: Container3D{
			Length: container.Length,
			Width:  container.Width,
			Height: container.Height,
		},
		Items: items,
		Stats: Stats3D{
			TotalItems:      len(placements),
			ContainerVolume: containerVolume,
			ItemsVolume:     totalItemVolume,
			UtilizationRate: utilizationRate,
		},
	}

	// 转换为JSON
	jsonData, err := json.MarshalIndent(vizData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// 写入文件
	return os.WriteFile(filename, jsonData, 0644)
}

// generateSpaceUtilizationStats 生成空间利用率统计
func generateSpaceUtilizationStats(placements []ItemPlacementRecord, container Container) {
	totalContainerVolume := int64(container.Length) * int64(container.Width) * int64(container.Height)
	var totalItemVolume int64

	for _, placement := range placements {
		itemVolume := int64(placement.Orientation.Length) *
			int64(placement.Orientation.Width) *
			int64(placement.Orientation.Height)
		totalItemVolume += itemVolume
	}

	utilizationRate := float64(totalItemVolume) / float64(totalContainerVolume) * 100

	fmt.Printf("======================== Space Utilization ========================\n")
	fmt.Printf("Container Volume: %d cubic units\n", totalContainerVolume)
	fmt.Printf("Items Volume: %d cubic units\n", totalItemVolume)
	fmt.Printf("Space Utilization: %.2f%%\n", utilizationRate)
	fmt.Printf("================================================================\n\n")
}

// generateHTMLViewer 提示用户使用HTML查看器
func generateHTMLViewer() error {
	fmt.Printf("HTML viewer is already available as 'bin_packing_viewer.html'\n")
	return nil
}
