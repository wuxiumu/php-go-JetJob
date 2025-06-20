<?php

return [

    'models' => [

        'permission' => Spatie\Permission\Models\Permission::class,
        'role' => Spatie\Permission\Models\Role::class,

    ],

    'table_names' => [

        'roles' => 'roles',
        'permissions' => 'permissions',
        'model_has_permissions' => 'model_has_permissions',
        'model_has_roles' => 'model_has_roles',
        'role_has_permissions' => 'role_has_permissions',

    ],

    'column_names' => [

        'model_morph_key' => 'model_id',

    ],

    'cache_expiration_time' => 60 * 24, // 24 hours

    'display_permission_in_exception' => false,

    'display_role_in_exception' => false,

    'register_permission_check_method' => true,

    'guard_name'=> 'api'

];
