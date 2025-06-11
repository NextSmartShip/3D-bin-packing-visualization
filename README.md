# 3D Bin Packing Visualization

A high-performance 3D bin packing algorithm implementation with interactive web-based visualization, written in Go with Three.js frontend.

![3D Bin Packing Demo](https://img.shields.io/badge/Go-1.24+-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Status](https://img.shields.io/badge/status-active-brightgreen.svg)

## 🚀 Features

### Core Algorithm
- **Advanced 3D Bin Packing**: Efficient algorithm for optimal space utilization
- **Multiple Packing Strategies**: Support for different item sorting methods
  - Volume-based sorting (default)
  - Base area priority sorting
  - Dimension-based sorting
- **Collision Detection**: Ensures no overlapping items
- **Multiple Orientations**: Automatically tests different item rotations
- **Time Management**: Configurable execution time limits with early exit strategies

### Interactive 3D Visualization
- **Real-time 3D Rendering**: Powered by Three.js
- **Interactive Controls**: 
  - Mouse drag to rotate view
  - Mouse wheel for zoom in/out
  - Middle mouse drag for panning
  - Touch support for mobile devices
- **Visual Features**:
  - Color-coded items with transparency
  - Wireframe toggle
  - Coordinate axes display
  - Container outline visualization
- **Responsive UI**: Expandable/collapsible items list for better UX

### Statistics & Analytics
- **Space Utilization**: Real-time calculation of packing efficiency
- **Volume Analysis**: Container vs. items volume comparison
- **Item Details**: Position coordinates and dimensions for each item
- **Performance Metrics**: Execution time and optimization statistics

## 📋 Requirements

- **Go**: Version 1.24 or higher
- **Modern Web Browser**: Chrome, Firefox, Safari, or Edge with WebGL support
- **Operating System**: Linux, macOS, or Windows

## 🛠️ Installation

1. **Clone the repository**:
```bash
git clone <repository-url>
cd 3D-bin-packing-visualization
```

2. **Initialize Go module** (if not already done):
```bash
go mod init 3d-bin-packing-visualization
go mod tidy
```

3. **Build the project**:
```bash
go build .
```

## 🎯 Usage

### Basic Usage

1. **Run the application**:
```bash
go run .
```

2. **Open the visualization**:
   - The program generates `bin_packing_3d.json` with packing results
   - Open `bin_packing_viewer.html` in your web browser
   - Interact with the 3D visualization using mouse/touch controls

### Customizing Items

Edit the `main.go` file to define your own container and items:

```go
func main() {
    container := Container{
        Length: 600,
        Width:  400,
        Height: 400,
    }
    
    items := []*Item{
        {
            Length: 380,
            Width:  320,
            Height: 100,
            Qty:    2,  // Quantity of this item type
        },
        // Add more items...
    }
    
    result := CanPack(container, items)
    fmt.Println(result)
}
```

### Advanced Configuration

Use packing options for different strategies:

```go
// Use base area sorting strategy
result := CanPack(container, items, WithBaseAreaSortFunc())

// Use dimension-based sorting
result := CanPack(container, items, WithDimensionSortFunc())

// Set custom warning execution time
result := CanPack(container, items, WithWarningExecutionTime(5*time.Second))

// Combine multiple options
result := CanPack(container, items, 
    WithBaseAreaSortFunc(),
    WithWarningExecutionTime(10*time.Second))
```

## 🏗️ Project Structure

```
3D-bin-packing-visualization/
├── main.go                    # Application entry point with sample data
├── 3d_bin_packing.go         # Core packing algorithm implementation
├── bin_packing_viewer.html   # Interactive 3D visualization interface
├── go.mod                    # Go module dependencies
├── bin_packing_3d.json      # Generated packing results (runtime)
└── README.md                 # Project documentation
```

## 🔧 Technical Details

### Architecture
The project follows Domain-Driven Design (DDD) principles:

- **Domain Layer**: Pure business entities (`Item`, `Container`, `Position`)
- **Command Layer**: Business logic and packing algorithms
- **Interface Layer**: JSON generation and web interface
- **Infrastructure Layer**: File I/O and external dependencies

### Key Components

#### Data Structures
- `Item`: Represents items to be packed with dimensions and quantity
- `Container`: Defines the container constraints
- `PlacedItem`: Represents successfully packed items with positions
- `PackingSession`: Manages the entire packing process state

#### Algorithms
- **Space-First Fit**: Finds optimal positions for items
- **Orientation Testing**: Tries all possible item rotations
- **Collision Detection**: Ensures spatial constraints
- **Backtracking**: Optimizes placement through recursive search

#### Visualization
- **Three.js Integration**: 3D rendering with WebGL
- **Responsive Design**: Mobile and desktop compatible
- **Interactive Controls**: Intuitive navigation and exploration
- **Real-time Updates**: Dynamic UI based on packing results

## 📊 Performance

- **Typical Performance**: Handles 10-50 items efficiently
- **Time Complexity**: Optimized with early exit strategies
- **Memory Usage**: Minimal memory footprint
- **Scalability**: Configurable execution time limits for large datasets

## 🎮 Controls

### Mouse Controls
- **Left Click + Drag**: Rotate the 3D view
- **Middle Click + Drag**: Pan the camera
- **Mouse Wheel**: Zoom in/out
- **Right Click**: Disabled (no context menu)

### Touch Controls (Mobile)
- **Single Touch + Drag**: Rotate view
- **Pinch**: Zoom (if supported by browser)

### Interface Buttons
- **Reset View**: Return to default camera position
- **Toggle Wireframe**: Show/hide item edges
- **Toggle Axes**: Show/hide coordinate axes
- **Expand/Collapse**: Show all items in the list

## 🤝 Contributing

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Development Guidelines
- Follow Go best practices and formatting (`go fmt`)
- Add tests for new functionality
- Update documentation for API changes
- Ensure cross-platform compatibility

## 📈 Future Enhancements

- [ ] REST API for remote item submission
- [ ] Multiple container support
- [ ] Weight constraints and load balancing
- [ ] Export functionality (STL, OBJ formats)
- [ ] Batch processing capabilities
- [ ] Performance benchmarking tools
- [ ] Docker containerization

## 🐛 Troubleshooting

### Common Issues

**"undefined: Container" Error**:
- Make sure to run `go run .` instead of `go run main.go`
- Ensure both `main.go` and `3d_bin_packing.go` are in the same directory

**Visualization Not Loading**:
- Check browser console for JavaScript errors
- Ensure `bin_packing_3d.json` exists in the same directory as the HTML file
- Verify internet connection for Three.js CDN

**Performance Issues**:
- Reduce the number of items or container size
- Adjust `WarningExecutionTime` for faster results
- Use simpler sorting strategies for quick tests

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Three.js**: 3D visualization library
- **Go Community**: For excellent tooling and libraries
- **WebGL**: For hardware-accelerated 3D graphics
- **Contributors**: All developers who have contributed to this project

---

**Made with ❤️ using Go and Three.js**

For questions, issues, or contributions, please open an issue on the repository. 