function getData(){
    let nama = document.getElementById("floatingName").value
    let email = document.getElementById("floatingEmail").value
    let phoneNumber = document.getElementById("floatingPhoneNumber").value
    let subject = document.getElementById("floatingSelect").value
    let message = document.getElementById("floatingYourMessage").value

    if(nama == "") {
        return alert("Nama Harus Diisi")
    } else if(email == "") {
        return alert("Email Harus Diisi")
    } else if(phoneNumber == "") {
        return alert("No Telpon Harus Diisi")
    } else if(subject == "") {
        return alert("Kepentingan Harus Diisi")
    } else if(message == "") {
        return alert("Pesan Harus Diisi")
    }

    const emailReceiver= "yogaaneh1@gmail.com"

    let a = document.createElement("a")
    a.href = `mailto:${emailReceiver}?subject=${subject}&body=Halo,%0D%0ANama saya adalah ${nama}, bisakah menghubungi saya di ${phoneNumber},%0D%0ARincian kepentingan saya adalah ${message}.%0D%0A%0D%0ATerima Kasih,%0D%0A%0D%0ARegards,%0D%0A${nama}`
    
    a.click()

    // let data = {
    //     nama,
    //     email,
    //     phoneNumber,
    //     subject,
    //     message
    // }

    // console.log(data)
    
    let url = location.origin

    location.href = url + "/contact-me.html";

}