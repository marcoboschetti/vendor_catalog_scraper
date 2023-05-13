var categories;

$(document).ready(function () {
    var prodID = getUrlParameter('it');
    $.get("/api/admin/products/" + prodID, function (data) {
        console.log(data)
        renderProduct(data.product)
    });
});

function renderProduct(product) {
$("#mainImg").attr("src", "/"+product.image_urls[0])
$("#mainTitle").html(product.page_title)
$("#mainInfoBlock").html(product.info_block)
$("#mainDescription").html(product.description.replaceAll("col-sm-12","col-md-6 col-sm-12"))
$("#mainDescription2").html(product.description_2)


var sizeHtml = '<span>Selecione una talla:</span>';
var isFirst = true;
product.sizes.split("|").forEach(s => {
    var c = "";
    if(isFirst){
        isFirst = false
        c = "active"
    }

    sizeHtml += `
    <label class="`+c+`" for="`+s+`">`+s+`
    <input type="radio" id="`+s+`">
    </label>`;
});
$("#sizeContainer").html(sizeHtml);

}

var getUrlParameter = function getUrlParameter(sParam) {
    var sPageURL = window.location.search.substring(1),
        sURLVariables = sPageURL.split('&'),
        sParameterName,
        i;

    for (i = 0; i < sURLVariables.length; i++) {
        sParameterName = sURLVariables[i].split('=');

        if (sParameterName[0] === sParam) {
            return sParameterName[1] === undefined ? true : decodeURIComponent(sParameterName[1]);
        }
    }
    return false;
};