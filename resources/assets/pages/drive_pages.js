async function downloadFile(file_name, id) {
    try {
        let progress = 0;

        const resp = await fetch(`/api/v1/files/download?file_id=${id}`, { credentials: "same-origin" })
        if (resp.status != 200) {
            throw new Error(`${(await resp.json()).message}`)
        }
        const file_length = parseInt(resp.headers.get('Content-Length') ?? '0');

        const received = [];
        const reader = resp.body?.getReader()
        while (true) {
            if (reader == undefined) break;
            const { done, value } = await reader.read();
            if (done) break;
            received.push(value);
            progress += value.length;
            download_progress = (progress / file_length * 100)
        }
        const blob = new Blob(received, { type: resp.headers.get("Content-Type") ?? undefined });

        const link = document.createElement('a');
        link.href = window.URL.createObjectURL(blob);
        link.setAttribute('download', `${file_name}`);
        document.body.appendChild(link);
        link.click();
        link.parentNode?.removeChild(link);
    } catch (error) {
        console.log(error)
    }
}

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

    if (files_elm == null) {return}
    try {
        let formData = new FormData()

        try {
            const directory = window.location.search.split("?dir=")[1].split("&")[0]
            formData.append("dir", directory)
        } catch (error) {
            formData.append("dir", "/")
        }
        for (let i=0;i<files_elm.files.length;i++){
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

function openFolder() {
    let dir = document.getElementById("dir").value
    if (dir != ""){
        if (dir[0] != "/"){
            dir = `/${dir}`
        }
        window.location.search = "?dir="+dir
    }else{
        window.location.search = ""
    }
}