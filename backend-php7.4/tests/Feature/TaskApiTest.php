<?php

namespace Tests\Feature;

use App\Models\Task;
use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Spatie\Permission\Models\Role;
use Tests\TestCase;
use Tymon\JWTAuth\Facades\JWTAuth;

class TaskApiTest extends TestCase
{
    use RefreshDatabase; // 每个测试自动回滚数据库

    protected $token;

    protected function setUp(): void
    {
        parent::setUp();

        // 确保角色 `user` 存在
        Role::create(['name' => 'user']);

        // 使用 UserFactory 创建用户并分配角色
        $user = User::factory()->create([
            'email' => 'user1@example.com',
            'password' => bcrypt('password123'),
        ]);
        $user->assignRole('user');

        // 获取 JWT 令牌
        try {
            $this->token = JWTAuth::fromUser($user);
        } catch (\Tymon\JWTAuth\Exceptions\JWTException $e) {
            $this->fail('Could not create token: ' . $e->getMessage());
        }
    }

    protected function getAuthenticatedUserToken()
    {
        return $this->token;
    }

    protected function loginUserAndGetToken()
    {
        return $this->token;
    }

    public function test_authenticated_user_can_get_task_list()
    {
        // 预插入任务
        Task::factory()->count(3)->create();

        $response = $this->withHeaders([
            'Authorization' => 'Bearer ' . $this->token,
        ])->get('/api/v1/tasks');

        $response->assertStatus(200)
            ->assertJsonStructure([
                '*' => [
                    'id',
                    'name',
                    'type',
                    'command',
                    'schedule',
                    'status',
                ],
            ]);
    }

    public function test_unauthorized_user_cannot_get_task_list()
    {
        $response = $this->get('/api/v1/tasks');
//        $response->assertStatus(200);
        $response->assertStatus(401); // 未登录
    }

    public function test_can_generate_token()
    {
        $token = $this->loginUserAndGetToken();
        $this->assertNotNull($token);
    }

    public function test_authorized_user_can_get_task_details()
    {
        // 预插入一个任务
        $task = Task::factory()->create();

        $response = $this->withHeaders([
            'Authorization' => 'Bearer ' . $this->token,
        ])->get('/api/v1/tasks/' . $task->id);

        $response->assertStatus(200)
            ->assertJsonStructure([
                'id',
                'name',
                'type',
                'command',
                'schedule',
                'status',
            ]);
    }
}
