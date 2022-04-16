function logout() {
$.ajax({
        type: "GET",
        url: "/v1/logout",
        success: function (data) {
            window.location.reload(true)
        }
    })
}