# Create plugins in Go

The plugin mechanism in Go provides a way to dynamically load shared objects, allowing
applications to extend their functionality at runtime without the need for recompilation.
This feature allows us to use the microkernel architecture style.

Plugins are created as shared libraries using the plugin package. 
A plugin is essentially a Go package compiled with the -buildmode=plugin flag, resulting in a .so
file on Unix-like systems. Windows is not supported.

Once compiled, these plugins can be dynamically loaded by a Go program using the plugin.Open
function, which returns a *plugin.Plugin object.
Symbols within the plugin, such as functions or variables, can then be
accessed using the Lookup method.

[Documentation of plugin's package](https://pkg.go.dev/plugin)
