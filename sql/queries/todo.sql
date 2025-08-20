-- name: CreateTodo :one
INSERT INTO todos (id,title,description,status,created_at,updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllTodos :many
SELECT * FROM todos;


-- name: MarkTodoAsDone :one
UPDATE todos 
SET status = $1,
updated_at =$2
WHERE id = $3
RETURNING *;

-- name: DeleteAtodo :exec
DELETE FROM todos
where id = $1;

-- name: FilterTodos :many
SELECT * FROM todos
WHERE status = $1;

-- name: GetTodoByID :one
SELECT * FROM todos
WHERE id = $1; 