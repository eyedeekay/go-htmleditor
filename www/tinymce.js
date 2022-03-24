tinymce.init({
    selector: '#editor',
    plugins: 'preview importcss searchreplace autolink directionality code visualblocks visualchars fullscreen image link media template codesample table charmap pagebreak nonbreaking anchor insertdatetime advlist lists wordcount help charmap quickbars emoticons',
    menubar: 'custom edit view insert format tools table help',
    toolbar: 'undo redo | bold italic underline strikethrough | fontselect fontsizeselect formatselect | alignleft aligncenter alignright alignjustify | outdent indent |  numlist bullist | forecolor backcolor removeformat | pagebreak | charmap emoticons | fullscreen preview | insertfile image media template link anchor codesample | ltr rtl',
    menu: {
        custom: {
            title: 'File',
            items: 'newfile preview savefile saveaswhat loadfile'
        }
    },
    setup: function(editor) {
        editor.ui.registry.addMenuItem('newfile', {
            text: 'New File',
            icon: 'new-document',
            context: 'file',
            onAction: function() {
                // create a dialog to enter the name of the new file
                tinymce.activeEditor.windowManager.open({
                    title: 'New File',
                    body: {
                        type: 'panel',
                        items: [{
                            type: 'htmlpanel',
                            html: 'Enter Filename:'
                        }, {
                            type: 'input',
                            name: 'filename',
                            placeholder: 'Enter Filename'
                        }]
                    },
                    buttons: [{
                        type: 'cancel',
                        text: 'Cancel'
                    }, {
                        type: 'submit',
                        text: 'OK'
                    }],
                    onSubmit: function(e) {
                        var data = e.getData();
                        var xhr = new XMLHttpRequest();
                        xhr.open("POST", "/new", true);
                        xhr.setRequestHeader("Content-Type", "application/json");
                        xhr.send(JSON.stringify(data));
                        xhr.onreadystatechange = function() {
                            if (xhr.readyState === 4 && xhr.status === 200) {
                                //close the dialog
                                document.getElementById("filename").textContent = data.filename;
                                tinymce.activeEditor.windowManager.close();
                                editor.setContent(xhr.responseText);
                            }
                        }
                    }
                });
            }
        });
        editor.ui.registry.addMenuItem('preview', {
            text: 'Preview',
            icon: 'preview',
            context: 'file',
            onAction: function() {
                tinymce.activeEditor.execCommand('mcePreview');
            }
        });


        editor.ui.registry.addMenuItem('savefile', {
            text: 'Save',
            icon: 'save',
            context: 'file',
            onAction: function() {
                //gather everything in the textarea and send it to the server /save
                var text = editor.getContent();
                var data = {
                    "path": "tinymce.html",
                    "text": text,
                };
                var xhr = new XMLHttpRequest();
                xhr.open("POST", "/save", true);
                xhr.setRequestHeader("Content-Type", "application/json");
                xhr.send(JSON.stringify(data));
            }
        });

        editor.ui.registry.addNestedMenuItem('saveaswhat', {
            text: 'Save As',
            icon: 'save',
            getSubmenuItems: function() {
                return [{
                    type: 'menuitem',
                    text: 'Download file in browser',
                    onAction: function() {
                        var text = editor.getContent();
                        var data = {
                            "path": "tinymce.html",
                            "text": text,
                        };
                        var xhr = new XMLHttpRequest();
                        xhr.open("POST", "/download", true);
                        xhr.setRequestHeader("Content-Type", "application/json");
                        xhr.send(JSON.stringify(data));
                        xhr.onreadystatechange = function() {
                            if (this.readyState == 4 && this.status == 200) {
                                var blob = new Blob([this.response], {
                                    type: "text/plain;charset=utf-8"
                                });
                                var a = window.document.createElement('a');
                                a.href = window.URL.createObjectURL(blob);
                                a.download = data.path;
                                document.body.appendChild(a)
                                a.click();
                                document.body.removeChild(a)
                            }
                        }
                    }
                }];
            }
        });

        editor.ui.registry.addMenuItem('loadfile', {
            text: 'Load',
            icon: 'load',
            onAction: function() {
                var input = document.createElement('input');
                input.type = 'file';
                input.click();
                input.onchange = e => {
                    var file = e.target.files[0];
                    console.log(file);
                    var data = {
                        "path": file.name,
                    };
                    var xhr = new XMLHttpRequest();
                    xhr.open("POST", "/load", true);
                    xhr.setRequestHeader("Content-Type", "application/json");
                    xhr.send(JSON.stringify(data));
                    xhr.onreadystatechange = function() {
                        if (this.readyState == 4 && this.status == 200) {
                            console.log("loading", data.path)
                            editor.setContent(xhr.responseText);
                            document.getElementById("filename").textContent = data.path;
                        }
                    };
                }
            }
        });
    },
});