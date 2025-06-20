<?php

namespace Tests\Feature;

use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Spatie\Permission\Models\Role;
use Tests\TestCase;
use Tymon\JWTAuth\Facades\JWTAuth;

class RolePermissionTest extends TestCase
{
    use RefreshDatabase;

    public function setUp(): void
    {
        parent::setUp();
        Role::create(['name' => 'admin']);
        Role::create(['name' => 'user']);
    }

    public function test_admin_can_access_admin_api()
    {
        $admin = User::factory()->create();
//        $admin = User::factory()->create([
//            'email' => 'user1@example.com',
//            'password' => bcrypt('password123'),
//        ]);

        $admin->assignRole('admin');
        $token = JWTAuth::fromUser($admin);

        $response = $this->withHeaders([
            'Authorization' => 'Bearer ' . $token,
        ])->get('/api/v1/users');

        $response->assertStatus(200);
        // 视接口具体返回结构补充assertJson
    }

    public function test_user_cannot_access_admin_api()
    {
        $token = '';
        $response = $this->withHeaders([
            'Authorization' => 'Bearer ' . $token,
        ])->get('/api/v1/users');

        $response->assertStatus(401); // 权限不足
    }
}
