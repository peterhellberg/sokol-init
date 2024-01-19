const std = @import("std");

pub fn build(b: *std.Build) !void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    const sokol = b.dependency("sokol", .{
        .target = target,
        .optimize = optimize,
    });

    const exe = b.addExecutable(.{
        .name = "sokol-init",
        .target = target,
        .optimize = optimize,
        .root_source_file = .{
            .path = "src/main.zig",
        },
    });

    exe.root_module.addImport("sokol", sokol.module("sokol"));

    b.installArtifact(exe);

    b.step("run", "Run").dependOn(
        &b.addRunArtifact(exe).step,
    );
}
