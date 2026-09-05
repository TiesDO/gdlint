[
    (variable_statement
        "var" @var.var
        name: (_) @var.name
        "=" @var.assign
        value: (_) @var.value
    )
    (variable_statement
        "var" @var.var
        name: (_) @var.name
        ":" @var.colon
        type: (_) @var.type
        (
            "=" @var.assign
            value: (_) @var.value
        )?
    )
    (variable_statement
        "var" @var.var
        name: (_) @var.name
        (inferred_type) @var.infer
        value: (_) @var.value
    )
]
