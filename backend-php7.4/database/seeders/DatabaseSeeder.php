<?php

namespace Database\Seeders;

use Illuminate\Database\Seeder;

class DatabaseSeeder extends Seeder
{
    /**
     * Seed the application's database.
     *
     * @return void
     */
    public function run()
    {
        // \App\Models\User::factory(10)->create();
        // 调用 RolesTableSeeder
        $this->call(RolesTableSeeder::class);
        // 你可以在这里调用其他的 seeder 文件
    }
}
