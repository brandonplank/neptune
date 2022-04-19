const successNotification = window.createNotification({})
const errorNotification = window.createNotification({
    theme: 'error'
})

$("#passwordChange").submit(function(e) {
    e.preventDefault();
});

function changePassword() {
    var currentPassword = document.getElementById('passwordChange').elements['CurrentPassword'].value;
    var newPassword = document.getElementById('passwordChange').elements['NewPassword'].value;
    $.ajax({
        type: "POST",
        url: "/v1/changePassword",
        headers: {
            "Content-Type": "application/json"
        },
        data: JSON.stringify({
            "password": currentPassword,
            "newPassword": newPassword
        }),
        success: function (data) {
            successNotification({
                title: 'Success',
                message: data.message
            })
            setTimeout(function () {
                window.location.reload(true)
            }, 2*1000)
        },
        error: function (error) {
            try {
                var j = JSON.parse(error.responseText)
                errorNotification({
                    title: 'Error',
                    message: j.message
                })
            } catch (error) {
                console.error(error)
                errorNotification({
                    title: 'Error',
                    message: "An unknown error occurred"
                })
            }
        }
    })
}

$("#login").submit(function(e) {
    e.preventDefault();
});

function login() {
    var email = document.getElementById('login').elements['Email'].value;
    var password = document.getElementById('login').elements['Password'].value;
    console.log("Logging in")
    $.ajax({
        type: "POST",
        url: "/v1/login",
        headers: {
            "Content-Type": "application/json"
        },
        data: JSON.stringify({
            "email": email,
            "password": password
        }),
        success: function (data) {
            successNotification({
                title: 'Success',
                message: data.message
            })
            setTimeout(function () {
                window.location.reload(true)
            }, 2*1000)
        },
        error: function (error) {
            try {
                var j = JSON.parse(error.responseText)
                errorNotification({
                    title: 'Error',
                    message: j.message
                })
            } catch (error) {
                console.error(error)
                errorNotification({
                    title: 'Error',
                    message: "An unknown error occurred"
                })
            }
        }
    })
}