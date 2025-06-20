<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;

class RefreshToken
{
    /**
     * Handle an incoming request.
     *
     * @param \Illuminate\Http\Request $request
     * @param \Closure(\Illuminate\Http\Request): (\Illuminate\Http\Response|\Illuminate\Http\RedirectResponse) $next
     * @return \Illuminate\Http\Response|\Illuminate\Http\RedirectResponse
     */
    public function handle($request, Closure $next)
    {
        try {
            $response = $next($request);

            // 检查 Token 是否即将过期

            return $response;
        } catch (TokenExpiredException $e) {
            // Token 已过期，尝试刷新
            try {
                $newToken = auth()->refresh();
                return response()->json([
                    'message' => 'Token 已刷新',
                    'token' => $newToken
                ], 200);
            } catch (JWTException $e) {
                // 刷新失败，强制用户重新登录
                return response()->json([
                    'message' => '认证失败，请重新登录',
                    'code' => 401
                ], 401);
            }
        }
    }
}
