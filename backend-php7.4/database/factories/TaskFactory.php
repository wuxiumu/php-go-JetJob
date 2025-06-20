<?php

namespace Database\Factories;

use App\Models\Task;
use Illuminate\Database\Eloquent\Factories\Factory;

class TaskFactory extends Factory
{
    protected $model = Task::class;
    public function definition()
    {
        return [
            'name' => $this->faker->word,
            'type' => 'shell',
            'command' => 'echo "Hello World"',
            'schedule' => '* * * * *',
            'status' => 'active',
        ];
    }
}
