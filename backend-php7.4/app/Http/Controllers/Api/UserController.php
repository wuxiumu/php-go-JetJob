<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\Request;
use App\Models\User;
use Illuminate\Support\Facades\Hash;
use Spatie\Permission\Models\Role;
use Illuminate\Validation\Rule;
use Spatie\Permission\Exceptions\RoleDoesNotExist;

class UserController extends Controller
{
    // 列表
    public function index()
    {
        $users = User::with('roles', 'permissions')->paginate(20);
        return response()->json($users);
    }

    // 详情
    public function show($id)
    {
        $user = User::with('roles', 'permissions')->findOrFail($id);
        return response()->json($user);
    }

    // 新建
    public function store(Request $request)
    {
        $data = $request->validate([
            'name'     => 'required|string|max:255',
            'email'    => 'required|email|unique:users',
            'password' => 'required|string|min:6',
            'roles'    => 'array'
        ]);
        $user = User::create([
            'name'     => $data['name'],
            'email'    => $data['email'],
            'password' => Hash::make($data['password']),
        ]);
        // 分配角色，默认user
        $roles = $data['roles'] ?? ['user'];
        $user->assignRole($roles);
        return response()->json(['message' => '用户创建成功', 'user' => $user->load('roles')], 201);
    }

    // 更新
    public function update(Request $request, $id)
    {
        $user = User::findOrFail($id);
        $data = $request->validate([
            'name'     => 'sometimes|string|max:255',
            'email'    => ['sometimes','email', Rule::unique('users')->ignore($id)],
            'password' => 'sometimes|string|min:6',
            'roles'    => 'array'
        ]);
        if (isset($data['name']))     $user->name = $data['name'];
        if (isset($data['email']))    $user->email = $data['email'];
        if (isset($data['password'])) $user->password = Hash::make($data['password']);
        $user->save();
        if (isset($data['roles'])) {
            $user->syncRoles($data['roles']);
        }
        return response()->json(['message' => '用户信息已更新', 'user' => $user->load('roles')]);
    }

    // 删除
    public function destroy($id)
    {
        $user = User::findOrFail($id);
        $user->delete();
        return response()->json(['message' => '用户已删除']);
    }

    // 分配角色
    public function assignRoles(Request $request, $id)
    {
        $user = User::findOrFail($id);
        $roles = $request->input('roles', []);
        $res = [];
        foreach ($roles as $role){
            $res[] = Role::findOrCreate($role);
        }
        if (!is_array($roles) || empty($roles)) {
            return response()->json(['message' => '角色参数必须为非空数组', 'roles' => $roles], 422);
        }

        // 检查所有角色名是否存在
        $rolesInDb = Role::whereIn('name', $roles)->pluck('name')->toArray();
        // 表Role
        $missingRoles = array_diff($roles, $rolesInDb);

        if (!empty($missingRoles)) {
            \Log::warning('分配角色失败，数据库未找到角色', ['missing' => $missingRoles, 'input' => $roles]);
            return response()->json([
                'message' => '以下角色不存在',
                'roles'=>$roles,
                'missing_roles' => $missingRoles,
                'all_roles' => $rolesInDb
            ], 422);
        }

        try {
            $user->syncRoles($roles);
        } catch (RoleDoesNotExist $e) {
            \Log::error('RoleDoesNotExist异常', ['exception' => $e->getMessage(), 'roles' => $roles]);
            return response()->json(['message' => '指定的角色不存在', 'roles' => $roles,'res'=>$res], 400);
        } catch (\Exception $e) {
            \Log::error('分配角色其他异常', ['exception' => $e->getMessage(), 'roles' => $roles]);
            return response()->json(['message' => '分配角色失败，系统异常', 'roles' => $roles], 500);
        }

        return response()->json(['message' => '角色分配成功', 'user' => $user->load('roles')]);
    }

    // 分配权限
    public function assignPermissions(Request $request, $id)
    {
        $user = User::findOrFail($id);
        $permissions = $request->input('permissions', []);
        $user->syncPermissions($permissions);
        return response()->json(['message' => '权限分配成功', 'user' => $user->load('permissions')]);
    }
}
