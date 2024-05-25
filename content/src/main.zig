const std = @import("std");
const sokol = @import("sokol");

const app = sokol.app;
const gfx = sokol.gfx;
const log = sokol.log;
const glue = sokol.glue;

var pass: gfx.PassAction = .{};

pub fn main() void {
    app.run(.{
        .init_cb = init,
        .frame_cb = frame,
        .event_cb = input,
        .cleanup_cb = cleanup,
        .width = 1280,
        .height = 720,
        .icon = .{ .sokol_default = true },
        .window_title = "sokol-init",
        .logger = .{ .func = log.func },
        .win32_console_attach = true,
    });
}

export fn init() void {
    gfx.setup(.{
        .environment = glue.environment(),
        .logger = .{ .func = log.func },
    });

    pass.colors[0] = .{
        .load_action = .CLEAR,
        .clear_value = .{ .r = 0, .g = 0, .b = 0, .a = 1 },
    };

    std.debug.print("Backend: {}\n", .{gfx.queryBackend()});
}

export fn frame() void {
    const g = pass.colors[0].clear_value.g + 0.005;

    pass.colors[0].clear_value.g = if (g > 1.0) 0.3 else g;

    gfx.beginPass(.{ .action = pass, .swapchain = glue.swapchain() });
    gfx.endPass();
    gfx.commit();
}

export fn input(event: ?*const app.Event) void {
    const ev = event.?;

    if (ev.type == .KEY_DOWN) {
        switch (ev.key_code) {
            .Q, .ESCAPE => app.requestQuit(),
            .F => app.toggleFullscreen(),
            else => {},
        }
    }
}

export fn cleanup() void {
    gfx.shutdown();
}
