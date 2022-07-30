function logout() {
$.ajax({
        type: "GET",
        url: "/v1/logout",
        success: function (data) {
            window.location.reload(true)
        }
    })
}

var upload = document.getElementById("pdfs")
const reader = new FileReader()
var ReturnedBytes = new Uint8Array(1)

upload.onchange = function(e) {
    reader.readAsArrayBuffer(e.target.files[0]);
    reader.onloadend = (evt) => {
        if (evt.target.readyState === FileReader.DONE) {
            const arrayBuffer = evt.target.result
            const array = new Uint8ClampedArray(arrayBuffer)
            const bytes = ProcessPDF(array, array.length)
            const blob = new Blob([bytes], {type: "application/zip"})

            const link = document.createElement('a')
            link.href = window.URL.createObjectURL(blob)
            link.download = "codes.zip"
            link.click()
        }
    }
};

function GetWasmFile() {
    var bytes = Array.from(ReturnedBytes)
    console.log(bytes)
}