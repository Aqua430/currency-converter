document.getElementById("converter_form").addEventListener("submit", checkForm);
function checkForm(event){
    event.preventDefault();
    const form = document.getElementById('converter_form');
    const amount = form.elements["amount"].value;
    const from_currency = form.elements["fromCurrency"].value;
    const to_currency = form.elements["toCurrency"].value;
    let fail = ""
    if(amount == "")
        fail = "Ошибка: введите сумму";
    else if(from_currency === "Валюта" || to_currency === "Валюта")
        fail = "Ошибка: выберите валюту";

    const errorElement = document.getElementById('error');

    if(fail != "")
        errorElement.innerHTML = fail;
    else {
        errorElement.innerHTML = "";
        convert(amount, from_currency, to_currency);
    }
}

function convert(amount, from, to) {
    fetch(`https://currency-converter-production-2995.up.railway.app//convert?from=${from}&to=${to}&amount=${amount}`)
        .then(response => response.json())
        .then(data => {
            document.getElementById("result").textContent = data.result.toFixed(3);
        })
        .catch(error => {
            document.getElementById("error").innerHTML = "Ошибка при запросе к серверу";
            console.error(error);
        });
}
