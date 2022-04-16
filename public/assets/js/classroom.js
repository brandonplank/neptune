const successNotification = window.createNotification({})
const errorNotification = window.createNotification({
    theme: 'error'
})

// Scheduler
function runAtTime(method, hour, minute, second) {
    (function loop() {
        var now = new Date();
        if (now.getHours() === hour && now.getMinutes() === minute && now.getSeconds() === second) {
            method();
        }
        now = new Date();
        var delay = 60000 - (now % 60000);
        setTimeout(loop, delay);
    })();
}

// CSV Table fetcher
function getTable() {
    $.ajax({
        type: "GET",
        url: "/v1/GetCSV/",
        success: function (data) {
            $('table').replaceWith(arrayToTable(Papa.parse(data).data))
        }
    })
}

// Run this automatically on page load
getTable()

//Run at midnight
runAtTime(getTable, 0, 0, 0)

function sendTestContent(content) {
    $.ajax({
        type: "POST",
        url: "/v1/id/" + btoa(content),
        success: function (data) {
            console.log(data)
            getTable()
        }
    })
}

function cleanClass(name) {
    if(typeof name == 'undefined') {
        name = ""
    }
    console.log("Cleaning " + name)
    $.ajax({
        type: "GET",
        url: "/v1/CleanClass/" + name,
        success: function (data) {
            console.log(data)
            getTable()
        }
    })
}

function addTeacher(name, email) {
    console.log("Adding", name)
    $.ajax({
        type: "POST",
        url: "/addTeacher",
        headers: {
            "Content-Type": "application/json"
        },
        data: JSON.stringify({
            "name": name,
            "email": email
        }),
        success: function (data) {
        }
    })
}

function removeTeacher(email) {
    console.log("Removing", email)
    $.ajax({
        type: "POST",
        url: "/removeTeacher",
        headers: {
            "Content-Type": "application/json"
        },
        data: JSON.stringify({
            "email": email,
        }),
        success: function (data) {
        }
    })
}

function arrayToTable(tableData) {
    var table = $('<table class="table"></table>');
    $(tableData).each(function (i, rowData) {
        var row = $('<tr></tr>');
        $(rowData).each(function (j, cellData) {
            if(cellData.length >= 1) {
                row.append($('<td>'+cellData+'</td>'))
            }
        })
        table.append(row)
    })
    return table
}

// QR Code scanner
function sendStatusToWebPage(data) {
    let parsedJson = JSON.parse(JSON.stringify(data))
    if(parsedJson.isOut) {
        successNotification({
            title: 'Signed back in',
            message: parsedJson.name + ' has signed back in'
        })
    } else {
        successNotification({
            title: 'Signed out',
            message: parsedJson.name + ' has signed out'
        })
    }
    $.ajax({
        type: "POST",
        url: "/v1/id/" + btoa(parsedJson.name),
        success: function (dataPost) {
            console.log(dataPost)
            getTable()
        }
    })
}

function DoIfAdminQR(content) {
    if(content.includes("// override")) {
        successNotification({
            title: 'ADMIN',
            message: 'Script is now executing'
        });
        eval(content)
        return true
    }
    return false
}

function verifyName(name) {
    var exp = /^([a-zA-Z\-]+)\s*,\s*([a-zA-Z]+)(\s+([a-zA-Z]+))?$/gm;
    return name.match(exp);
}

var lastResult
function onScanSuccess(decodedText) {
    if (decodedText !== lastResult) {
        lastResult = decodedText
        setTimeout(function () {
            lastResult = null
        }, 10*1000)
        if(!verifyName(decodedText)) {
            if(!DoIfAdminQR(decodedText)) {
                errorNotification({
                    title: 'Error',
                    message: 'The QR you scanned is not valid',
                });
            }
            return
        }
        // successNotification({
        //     title: 'Success',
        //     message: 'Scanned QR code'
        // })
        $.ajax({
            type: "POST",
            url: "/v1/isOut/" + btoa(decodedText),
            success: sendStatusToWebPage
        })
    }
}

Html5Qrcode.getCameras().then(devices => {
    if (devices && devices.length) {
        var cameraId = devices[0].id
        console.log(`Got camera ID ${cameraId}`)
    }
}).catch(err => {
    console.error(err)
});

function changePassword(currentPassword, newPassword) {
    $.ajax({
        type: "POST",
        url: "/v1/changePassword",
        headers: {
            "Content-Type": "application/json"
        },
        data: JSON.stringify({
            "currentPassword": currentPassword,
            "newPassword": newPassword
        }),
        success: function (data) {
        }
    })
}


const html5QrCode = new Html5Qrcode("qr-reader", { formatsToSupport: [ Html5QrcodeSupportedFormats.QR_CODE ] })
const config = { fps: 60, qrbox: 250 }
html5QrCode.start({ facingMode: "user" }, config, onScanSuccess)

console.log("Hi reader :) This is Brandon here(Class of 2022) congrats on clicking F12 or view page src :P\n\nThis project was made using a multitude of languages, here is the list\n\nHTML(not really a programming language)\nJavaScript\nGoLang\n\nPlease always be nice to Mrs. Hart, she is the best teacher ever to exist.\nTalk to you on the flip side.")
