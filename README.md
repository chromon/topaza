# Topaza

## What is Topaza?

A lightweight concurrent server framework based on Golang

## Quick start

```
func main() {
    // server
	s := znet.NewServer()

	// router
	s.AddRouter(0, &HelloRouter{})

	// start
	s.Serve()
}
```

## License

The GPL3.0 License.

## Contact

Enjoy it.