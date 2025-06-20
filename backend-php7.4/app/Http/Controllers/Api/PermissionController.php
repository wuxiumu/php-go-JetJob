<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\Request;
use Spatie\Permission\Models\Permission;

class PermissionController extends Controller
{
    public function index()
    {
        return response()->json(Permission::all());
    }

    // 详情
    public function show($id)
    {
        $Permission = Permission::with('Permissions', 'permissions')->findOrFail($id);
        return response()->json($Permission);
    }

    public function store(Request $request)
    {
        $data = $request->validate(['name' => 'required|unique:Permissions']);
        $Permission = Permission::create(['name' => $data['name']]);
        return response()->json($Permission, 201);
    }

    public function update(Request $request, $id)
    {
        $Permission = Permission::findOrFail($id);
        $data = $request->validate(['name' => 'required|unique:Permissions,name,'.$id]);
        $Permission->name = $data['name'];
        $Permission->save();
        return response()->json($Permission);
    }

    public function destroy($id)
    {
        Permission::destroy($id);
        return response()->json(['message' => '角色已删除']);
    }
}
