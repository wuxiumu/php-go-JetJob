<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\Request;
use Spatie\Permission\Models\Role;

class RoleController extends Controller
{
    public function index()
    {
        return response()->json(Role::all());
    }

    // 详情
    public function show($id)
    {
        $role = Role::with('roles', 'permissions')->findOrFail($id);
        return response()->json($role);
    }

    public function store(Request $request)
    {
        $data = $request->validate(['name' => 'required|unique:roles']);
        $role = Role::create(['name' => $data['name']]);
        return response()->json($role, 201);
    }

    public function update(Request $request, $id)
    {
        $role = Role::findOrFail($id);
        $data = $request->validate(['name' => 'required|unique:roles,name,'.$id]);
        $role->name = $data['name'];
        $role->save();
        return response()->json($role);
    }

    public function destroy($id)
    {
        Role::destroy($id);
        return response()->json(['message' => '角色已删除']);
    }
}
