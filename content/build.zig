const std = @import("std");

pub fn build(b: *std.Build) !void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    const exe = b.addExecutable(.{
        .name = "sokol-init",
        .root_source_file = .{
            .path = "src/main.zig",
        },
        .target = target,
        .optimize = optimize,
    });

    const sokol = b.dependency("sokol", .{
        .target = target,
        .optimize = optimize,
    });

    exe.addModule("sokol", sokol.module("sokol"));
    exe.linkLibrary(sokol.artifact("sokol"));

    b.installArtifact(exe);

    b.step("run", "Run").dependOn(
        &b.addRunArtifact(exe).step,
    );
}
