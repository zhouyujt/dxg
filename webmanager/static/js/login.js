function load(showError){
    if(showError){
         alert("账号或密码错误");
    }
}

function checkInput() {
    var uname = document.getElementById("uname");
    if (uname.value.length == 0) {
        alert("请输入账号");
        return false;
    }

    var pwd = document.getElementById("pwd");
    if (pwd.value.length == 0) {
        alert("请输入密码");
        return false;
    }

    document.getElementById("pwdMD5").value = encodeMD5(pwd.value);
    return true;
}

function backHome() {
    window.location = "/";
    return false;
}

function encodeMD5(str){
    return $.md5(str);
}