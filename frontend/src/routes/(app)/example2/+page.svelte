<script>
    import {Login} from "$lib/requests/core/login"

    const login = Login.fromJsonString(`{"username":"first", "password":"password01"}`);
</script>

<div>
    <span>{login.username}</span>
    <span>{login.password}</span>
</div>