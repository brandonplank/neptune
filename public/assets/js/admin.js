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
        url: "/v1/GetAdminCSV/",
        success: function (data) {
            $('table').replaceWith(arrayToTable(Papa.parse(data).data))
        }
    })
}

function SearchStudent() {
    var name = document.getElementById('name').value
    console.log("Searching for", name)
    document.getElementById('name').value = "";
    $.ajax({
        type: "POST",
        url: "/v1/search/" + btoa(name),
        success: function (data) {
            $('table').replaceWith(arrayToTable(Papa.parse(data).data))
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

$('#name').keypress(function(e){
    if (e.which === 13) {
        SearchStudent()
    }
});

// Run this automatically on page load
getTable()

// Run every 10 seconds
setInterval(getTable, 1000 * 30);

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

try {
    var schoolSelect = document.getElementById("School");
    if(schoolSelect == null) {
        throw ""
    }
    $.ajax({
        type: "GET",
        url: "/v1/getSchools/",
        success: function (data) {
            console.log(data)
            for (var school of data) {
                var option = document.createElement("option");
                option.text = school.name;
                option.value = school.id;
                schoolSelect.add(option);
            }
        }
    })
} catch (e) {
    console.log("did not find school selection, not on page.")
}

try {
    var levelSelect = document.getElementById("Level");
    if(levelSelect == null) {
        throw ""
    }
    $.ajax({
        type: "GET",
        url: "/v1/getUserPermissionLevel/",
        success: function (data) {
            console.log(data.level)
            for (var i = 0; i < data.level; i++) {
                var option = document.createElement("option");
                option.text = i;
                option.value = i;
                levelSelect.add(option);
            }
        }
    })
} catch (e) {
}