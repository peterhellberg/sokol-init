const std = @import("std");
const sokol = @import("sokol");

const app = sokol.app;
const gfx = sokol.gfx;
const log = sokol.log;

var pa: gfx.PassAction = .{};

pub fn main() void {
    app.run(.{
        .init_cb = init,
        .frame_cb = frame,
        .cleanup_cb = cleanup,
        .width = 1280,
        .height = 720,
        .icon = .{ .sokol_default = true },
        .logger = .{ .func = log.func },
        .window_title = "sokol-init",
    });
}

export fn init() void {
    gfx.setup(.{
        .context = sokol.app_gfx_glue.context(),
        .logger = .{ .func = log.func },
    });

    pa.colors[0] = .{
        .load_action = .CLEAR,
        .clear_value = .{ .r = 1, .g = 0, .b = 0, .a = 1 },
    };

    std.debug.print("Backend: {}\n", .{gfx.queryBackend()});
}

export fn frame() void {
    const g = pa.colors[0].clear_value.g + 0.001;

    pa.colors[0].clear_value.g = if (g > 1.0) 0.0 else g;

    gfx.beginDefaultPass(pa, app.width(), app.height());
    gfx.endPass();
    gfx.commit();
}

export fn cleanup() void {
    gfx.shutdown();
}
