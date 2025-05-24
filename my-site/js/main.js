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
    if(fail != "")
        document.getElementById('error').innerHTML = fail;
    else {
        convert(amount, from_currency, to_currency);
    }
}

function convert(amount, from, to) {
    let rates = {
        USD: 1,
        RUB: 90,
        KZT: 450
    };
    const result = (amount / rates[from]) * rates[to];
    document.getElementById("result").textContent = result.toFixed(3);
}