<?php

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;
use App\Http\Controllers\Api;

// 健康检查
Route::get('/ping', function () {
    return response()->json(['pong' => true]);
});

Route::prefix('v1')->group(function () {
    // 用户认证
    Route::post('/login',    [Api\AuthController::class, 'login'])->middleware('throttle:10,1')->name('login');
    Route::post('/register', [Api\AuthController::class, 'register']);
    Route::post('/logout',   [Api\AuthController::class, 'logout'])->middleware('auth:api');
    Route::get('/me',        [Api\AuthController::class, 'me'])->middleware('auth:api');

    // 需要认证的接口
    Route::middleware(['auth:api', 'refresh.token', 'throttle:60,1'])->group(function () {

        // =================== 用户管理、角色、权限（建议仅 admin 角色能访问）====================

        // 用户管理
        Route::get('/users',               [Api\UserController::class, 'index']);
        Route::post('/users',              [Api\UserController::class, 'store']);
        Route::get('/users/{id}',          [Api\UserController::class, 'show']);
        Route::put('/users/{id}',          [Api\UserController::class, 'update']);
        Route::delete('/users/{id}',       [Api\UserController::class, 'destroy']);

        // 用户角色/权限分配
        Route::post('/users/{id}/roles',        [Api\UserController::class, 'assignRoles']);
        Route::post('/users/{id}/permissions',  [Api\UserController::class, 'assignPermissions']);

        // 角色管理
        Route::get('/roles',               [Api\RoleController::class, 'index']);
        Route::post('/roles',              [Api\RoleController::class, 'store']);
        Route::put('/roles/{id}',          [Api\RoleController::class, 'update']);
        Route::delete('/roles/{id}',       [Api\RoleController::class, 'destroy']);

        // 权限管理
        Route::get('/permissions',         [Api\PermissionController::class, 'index']);
        Route::post('/permissions',        [Api\PermissionController::class, 'store']);
        Route::put('/permissions/{id}',    [Api\PermissionController::class, 'update']);
        Route::delete('/permissions/{id}', [Api\PermissionController::class, 'destroy']);

        // =================== 任务管理 ===================
        Route::get('/tasks',               [Api\TaskController::class, 'index']);
        Route::post('/tasks',              [Api\TaskController::class, 'store']);
        Route::get('/tasks/{id}',          [Api\TaskController::class, 'show']);
        Route::put('/tasks/{id}',          [Api\TaskController::class, 'update']);
        Route::delete('/tasks/{id}',       [Api\TaskController::class, 'destroy']);

        // =================== 节点管理 ===================
        Route::get('/nodes',               [Api\NodeController::class, 'index']);
        Route::post('/nodes',              [Api\NodeController::class, 'store']);
        Route::get('/nodes/{id}',          [Api\NodeController::class, 'show']);
        Route::put('/nodes/{id}',          [Api\NodeController::class, 'update']);
        Route::delete('/nodes/{id}',       [Api\NodeController::class, 'destroy']);

        // =================== 日志管理 ===================
        Route::get('/logs',                [Api\LogController::class, 'index']);
        Route::get('/logs/{id}',           [Api\LogController::class, 'show']);

        // 其它业务接口可继续补充...
    });

    // 公共、无需登录的接口（如文档、验证码等，可放这里）
});
