async function deleteFile(id) {
    const confirmed = window.confirm("Are you sure you want to delete this file?");
    if (confirmed) {
        try {
            await fetch(`/api/v1/files/delete`,
                {
                    method: "DELETE",
                    body: JSON.stringify({ file_id: id }),
                    credentials: "same-origin",
                    headers: { "Content-Type": "application/json" }
                }
            )
            window.location.reload()
        } catch (error) {
            console.log(error)
        }
    }
}

async function uploadMultipleFiles() {
    const files_elm = document.getElementById("files")

    if (files_elm == null) { return }
    try {
        let formData = new FormData()

        try {
            const directory = window.location.search.split("?dir=")[1].split("&")[0]
            formData.append("dir", directory)
        } catch (error) {
            formData.append("dir", "/")
        }
        for (let i = 0; i < files_elm.files.length; i++) {
            /** @type {File} */
            const file = files_elm.files[i]
            formData.append("files", file)
        }

        await fetch(`/api/v1/files/upload/batch`, {
            method: "POST",
            body: formData,
            credentials: "same-origin",
        })
        window.location.reload()
    } catch (error) {
        console.log(error)
    }
}

async function uploadFileViaWS() {
    const TRANSMISSION_END = "TRANSMISSION_END"
    const FILENAME_END = "FILENAME_END"
    const DIRECTORY_END = "DIRECTORY_END"

    /** @type {File} */
    let file = document.getElementById("files").files[0]
    let protocol = "wss://"
    if (window.location.protocol == "http:") {
        protocol = "ws://"
    }
    let dir = ""
    try {
        dir = window.location.search.split("?dir=")[1].split("&")[0]
    } catch (error) {
        dir = "/"
    }

    const ws = new WebSocket(protocol + document.location.host + "/ws/v1/files/upload")

    ws.addEventListener("open", async () => {
        ws.send(file.name)
        ws.send(dir)

        const stream = file.stream()

        for await (const chunk of stream) {
            ws.send(chunk)
            delete chunk
        }

        ws.send(unpack(TRANSMISSION_END))
    })

    ws.addEventListener("message", (event) => {
        console.log(event.data)
    })

    ws.addEventListener("error", (event) => {
        console.log("error:", event.data)
        ws.close()
    })
}

async function uploadMultipleFileViaWS() {
    const TRANSMISSION_END = "TRANSMISSION_END"
    const SPLIT_END = "SPLIT_END"

    let alert = document.getElementById("uploading-alert")
    alert.className = "alert alert-primary"

    /** @type {[]File} */
    let files = document.getElementById("files").files
    let protocol = "wss://"
    if (window.location.protocol == "http:") {
        protocol = "ws://"
    }
    let dir = ""
    try {
        dir = window.location.search.split("?dir=")[1].split("&")[0]
    } catch (error) {
        dir = "/"
    }

    const ws = new WebSocket(protocol + document.location.host + "/ws/v1/files/upload/batch")

    let stop = 0
    ws.addEventListener("open", async () => {
        for (let i = 0; i < files.length; i++) {
            /** @type {File} */
            const file = files[i]

            ws.send(file.name)
            ws.send(dir)

            const stream = file.stream()
            let sent = 0

            for await (const chunk of stream) {
                if (stop) { break }

                ws.send(chunk)
                sent += chunk.byteLength
                delete chunk

                alert.hidden = false
                alert.innerHTML = `Uploading ${file.name} ${((sent - ws.bufferedAmount) / file.size * 100).toFixed(2)} %`
            }

            ws.send(unpack(SPLIT_END))
        }
        ws.send(unpack(TRANSMISSION_END))
    })

    ws.addEventListener("message", (event) => {
        alert.hidden = false
        alert.innerHTML = event.data
        if (event.data == "[Success]") {
            alert.className = "alert alert-success"
        } else {
            alert.className = "alert alert-danger"
        }
        stop = 1
        setTimeout(() => { window.location.reload() }, 1000)
    })

    ws.addEventListener("error", () => {
        ws.close()
    })
}

/**
 * 
 * @param {string} str 
 */
function unpack(str) {
    const EOF = ' '
    let encoder = new TextEncoder()
    let bytes = encoder.encode(str + EOF)
    return bytes
}

function openFolder() {
    let dir = document.getElementById("dir").value
    if (dir == "") {
        window.location.search = ""
    } else {
        if (dir.includes('/')) {
            return document.getElementById("dir") = ""
        }
        if (dir[0] != "/") {
            dir = `/${dir}`
        }
        if (window.location.search.includes("?dir=")) {
            window.location.search += dir
        } else {
            window.location.search = "?dir=" + dir
        }
    }
}

function goBack() {
    if (window.location.search.includes("?dir=")) {
        let curr_dir = window.location.search.split("?dir=")[1]
        if (curr_dir.includes("&")) {
            curr_dir = curr_dir.split("&")[0]
        }

        const splitted_dir = curr_dir.split("/")
        splitted_dir.pop()
        const new_dir = splitted_dir.join("/")

        window.location.search = "?dir=" + new_dir
    }
}


function showModalImage(file_id, filename) {
    document.getElementById("modalPreviewTitle").innerHTML = filename

    // reset first
    document.getElementById("modalImageBody").src = ""

    document.getElementById("modalImageBody").src = `/api/v1/files/show?file_id=${file_id}`
    document.getElementById("modalImageBody").alt = filename
}