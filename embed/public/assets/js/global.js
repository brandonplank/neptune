function logout() {
$.ajax({
        type: "GET",
        url: "/v1/logout",
        success: function (data) {
            window.location.reload(true)
        }
    })
}

var upload = document.getElementById("pdfs");
const reader = new FileReader();
var ReturnedBytes = new Uint8Array(1);

upload.onchange = function(e) {
    reader.readAsArrayBuffer(e.target.files[0]);
    reader.onloadend = (evt) => {
        if (evt.target.readyState === FileReader.DONE) {
            const arrayBuffer = evt.target.result
            const array = new Uint8ClampedArray(arrayBuffer);
            ProcessPDF(array, array.length);
        }
    }
};

function GetWasmFile() {
    console.log(ReturnedBytes)
}