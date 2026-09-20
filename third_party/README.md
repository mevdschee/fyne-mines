# third_party

## glfw

A copy of `github.com/go-gl/glfw/v3.3/glfw`, wired up through a `replace` in
go.mod. It carries one patch, in `glfw/src/nsgl_context.m`.

GLFW always asks macOS for `NSOpenGLPFAAccelerated`, which is a hard
constraint, so on a host without a GPU no pixel format matches and window
creation fails with "NSGL: Failed to find a suitable pixel format". The patch
retries on that failure, swapping the accelerated constraint for an explicit
`kCGLRendererGenericFloatID`, which is Apple's CPU renderer. A machine with a
GPU never reaches the retry and behaves exactly as before. Setting
`GLFW_SOFTWARE_RENDERER` skips the first attempt and forces the CPU renderer.

See https://github.com/glfw/glfw/issues/2080, which is still open upstream.

When updating this copy, reapply the patch and change `upstreamTreeSHA` in
`glfw_tree_rebuild.go` to anything that is not the previous value. The go build
cache does not look at C sources outside the package directory, so without that
the patched file is silently ignored and you get an unpatched binary. See
https://github.com/go-gl/glfw/issues/269.
