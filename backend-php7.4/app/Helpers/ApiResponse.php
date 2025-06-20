<?php

namespace App\Helpers;

class ApiResponse
{
    public static function success($data = [], $msg = 'success', $code = 200)
    {
        return response()->json([
            'code' => $code,
            'msg'  => $msg,
            'data' => $data,
        ], $code);
    }

    public static function error($msg = 'error', $code = 400, $data = [])
    {
        return response()->json([
            'code' => $code,
            'msg'  => $msg,
            'data' => $data,
        ], $code);
    }

    // app/Helpers/ApiResponse.php
    function apiResponse($data = null, $message = null, $code = 200)
    {
        $response = [
            'success' => $code >= 200 && $code < 300,
            'message' => $message,
            'data' => $data,
        ];

        return response()->json($response, $code);
    }
}
