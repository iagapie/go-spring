fetch("/", {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'X-Requested-With': 'XMLHttpRequest',
        'X_SPRING_REQUEST_HANDLER': 'todo::OnFetchData',
        'X_SPRING_REQUEST_PARTIALS': 'new_todo/foo'
    },
    body: JSON.stringify({"foo": "boo"})
}).then(d => d.json()).then(v => console.log(v))

const formData = new FormData();
formData.append("_handler", "todo::OnFetchForm")
formData.append("title", "test title")

fetch("/", {
    method: "POST",
    body: formData,
    mode: 'cors'
}).then(d => d.json()).then(v => console.log(v))
