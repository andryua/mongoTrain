function funcSelect(sel) {
    var defgp = document.querySelectorAll(".default");
    var maingp = document.querySelectorAll(".main");
    var usrgp = document.querySelectorAll(".users");
    var keyarrdef = Object.keys(defgp);
    var keyarrmain = Object.keys(maingp);
    var keyarrusr = Object.keys(usrgp);
    switch (sel.value) {
        default:
            break;
        case "main":
            document.getElementById("collapsedit").hidden = true;
            break;
        case "default":
            document.getElementById("collapsedit").hidden = false;
            keyarrdef.forEach(function (key){
                defgp[key].hidden = true;
            });
            keyarrmain.forEach(function (key){
                maingp[key].hiddent = false;
            });
            keyarrusr.forEach(function (key){
                usrgp[key].hidden = true;
            });
            break;
        case "users":
            document.getElementById("collapsedit").hidden = false;
            keyarrdef.forEach(function (key){
                defgp[key].hidden = false;
            });
            keyarrmain.forEach(function (key){
                maingp[key].hidden = true;
            });
            keyarrusr.forEach(function (key){
                usrgp[key].hidden = true;
            });
            break;
    }
}
function funcType(type) {
    switch (type.class) {
        default:
            type.innerText = "none";
            break;
        case ".main":
            type.innerText = "головна";
            break;
        case ".default":
            type.innerText = "типова";
            break;
        case ".users":
            type.innerText = "користувацька";
            break;
    }
}

function funcSaveRule(id) {
    //console.log(id)
    $.post(
        "/update?id=" + id,
        $("#update-form-" + id).serialize(),
        function(data) {
            if (data == "saved") {
                $("#res-" + id).html("<div class='alert alert-success alert-dismissible fade show' role='alert'>Зміни збережено!<button type='button' class='btn-close' data-bs-dismiss='alert' aria-label='Close'></button></div>");
            } else {
                $("#res-" + id).html("<div class='alert alert-danger alert-dismissible fade show' role='alert'>Зміни не збережено!<button type='button' class='btn-close' data-bs-dismiss='alert' aria-label='Close'></button></div>");
            }
        },
    );
}


k=1
function funcAddManual(btn,id) {
    //console.log($(btn).parents());
    //var parentDiv = $(btn).parents()[1];
    //var id = btn.id.split("-")[1];
    var manual = document.getElementById("hasmanual-"+id);
    //btn.hidden = true;
    //document.getElementById("manualvn_text-"+k).value = document.getElementById("manualvn_text-"+k).textContent;
    k++;
    manual.innerHTML += "<div class=\"col-sm-4\"><input type=\"text\" name=\"manualValueName\" aria-label=\"ValueName\" class=\"form-control\" id=\"manualvn_text-" + k +"\" inputmode=\"text\" placeholder=\"ім'я значення\"></div>\n" +
        "                    <div class=\"col-sm-4\"><input type=\"text\" name=\"manualValue\" aria-label=\"Value\" class=\"form-control\" id=\"manualv_text-"+ k+"\" inputmode=\"text\" placeholder=\"значення\"></div>\n" +
        "                    <div class=\"col-sm-3\"><input type=\"text\" name=\"manualDescription\" aria-label=\"Info\" class=\"form-control\" id=\"manuali_text-"+k+"\" inputmode=\"text\" placeholder=\"опис\"></div>\n"
    //$($(btn).parents()[0]).find("input").attr("readonly", true);
}