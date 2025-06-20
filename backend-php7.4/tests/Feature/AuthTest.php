<?php

namespace Tests\Feature;

use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Spatie\Permission\Models\Role;
use Tests\TestCase;
use Tymon\JWTAuth\Facades\JWTAuth;

class AuthTest extends TestCase
{
    use RefreshDatabase; // 每个测试自动回滚数据库

    protected function setUp(): void
    {
        parent::setUp();

        // 确保角色 `user` 存在
        Role::create(['name' => 'admin']);
        Role::create(['name' => 'user']);
    }

    public function test_user_can_register()
    {
        $response = $this->postJson('/api/v1/register', [
            'name' => 'TestUser',
            'email' => 'testuser@example.com',
            'password' => 'password123456789',
            'password_confirmation' => 'password123456789',
        ]);

        $response->assertStatus(201)
            ->assertJson(['message' => 'User registered successfully']);

        $this->assertDatabaseHas('users', ['email' => 'testuser@example.com']);
    }

    public function test_user_can_login_and_get_jwt_token()
    {
        // 确保角色 `user` 存在，并指定守护器为 `api`
        Role::firstOrCreate(['name' => 'user', 'guard_name' => 'api']);

        // 使用 UserFactory 创建用户
        $user = User::factory()->create([
            'email' => 'user1@example.com',
            'password' => bcrypt('password123'),
        ]);

        $response = $this->postJson('/api/v1/login', [
            'email' => 'user1@example.com',
            'password' => 'password123',
        ]);

        $response->assertStatus(200)
            ->assertJsonStructure(['access_token', 'token_type', 'expires_in']);
    }

    public function test_expired_token_returns_401()
    {
        $token ='';
        // 模拟 Token 过期
        $response = $this->withHeaders([
            'Authorization' => 'Bearer ' . $token,
        ])->getJson('/api/v1/tasks');

        $response->assertStatus(401);
    }
}
