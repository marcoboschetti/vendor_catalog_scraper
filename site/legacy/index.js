var loadedCatalogProducts;
var loadedCatalogNewLabeledProducts;

var categoriesKeys = [];
var subcategoriesKeys = {};

var curCategoryIdx = 0;
var curSubcategoryIdx = 0;
var curProductUrlIdx = 0;

var totalProducts = 0;
var totalLabeledProducts = 0;
var processedProducts = 0;

$("#downloadOnlyNewProductImages").click(function () {
    loadedCatalogProducts = loadedCatalogOnlyNewProducts;
    startDownloadingProducts();
})

$("#downloadOnlyLabeledProductImages").click(function () {
    loadedCatalogProducts = loadedCatalogNewLabeledProducts;
    startDownloadingProducts();
})

$("#downloadAllProductImages").click(function () {
    startDownloadingProducts();
})

function startDownloadingProducts() {
    $(".downloadImagesBtnContainer").slideUp();
    $("#downloadProgressBarContainer").slideDown();

    $.each(loadedCatalogProducts, function (category, subcategories) {
        categoriesKeys.push(category);
        subcategoriesKeys[category] = [];
        $.each(subcategories, function (subcategory, productUrls) {
            subcategoriesKeys[category].push(subcategory);
            totalProducts += productUrls.length;
        });
    });

    curProductUrlIdx = -1;
    processNextItem();
}

function processNextItem() {
    var curProductUrl = nextCategoryIteration();
    ended = (curProductUrl == null);

    processedProducts += 1;
    updateimagesProgressBar(totalProducts, processedProducts);

    if (!ended) {
        loadProductFromURL(curProductUrl);
    } else {
        // Save processed records
        updateProcessedCategoriesRecord();
        $("#downloadProgressBarContainer").html("<h1>Finished!</h1>");
    }
}


function nextCategoryIteration() {
    var curCategoryKey = categoriesKeys[curCategoryIdx];
    var curSubcategory = subcategoriesKeys[curCategoryKey][curSubcategoryIdx];
    var currentProductUrlsArr = loadedCatalogProducts[curCategoryKey][curSubcategory];

    curProductUrlIdx += 1;
    if (!currentProductUrlsArr || curProductUrlIdx >= currentProductUrlsArr.length) {
        // Is at end of productsUrl, go to next subcategory
        curProductUrlIdx = 0;
        curSubcategoryIdx += 1;

        if (curSubcategoryIdx >= subcategoriesKeys[curCategoryKey].length) {
            // End of subcategory, go to next category
            curSubcategoryIdx = 0;
            curCategoryIdx += 1;

            if (curCategoryIdx >= categoriesKeys.length) {
                // ENDED!
                return null;
            }
        }
    }

    curCategoryKey = categoriesKeys[curCategoryIdx];
    curSubcategory = subcategoriesKeys[curCategoryKey][curSubcategoryIdx];
    currentProductUrlsArr = loadedCatalogProducts[curCategoryKey][curSubcategory];

    if (!loadedCatalogProducts[curCategoryKey][curSubcategory]) {
        return false;
    }

    var curProductUrl = loadedCatalogProducts[curCategoryKey][curSubcategory][curProductUrlIdx];
    return curProductUrl;
}


function loadProductFromURL(productUrl) {
    $.post("/api/legacy/get_product?", JSON.stringify({ product_url: productUrl }), function (data) {

        var curCategoryKey = categoriesKeys[curCategoryIdx];
        var curSubcategory = subcategoriesKeys[curCategoryKey][curSubcategoryIdx];

        // Write html into iframe
        $('#productsContainer').unbind("load");
        $('#productsContainer').on("load", function () {
            resizeIFrameToFitContent();
            
            setTimeout(function(){
                saveIframeScreenshot(curCategoryKey, curSubcategory, productUrl, processNextItem);
            }, 500);
        });

        var productHtml = htmlDecode(data.product_http);
        var doc = document.getElementById('productsContainer').contentWindow.document;

        $("#productsContainer").contents().find("html").html('');
        doc.open();
        doc.write(productHtml);
        doc.close();
    }, "json")
        .fail(function (response) {
            console.log('Failed to get product: ' + productUrl, "!", response);
            processNextItem();
        });;
}

function resizeIFrameToFitContent() {
    var iframe = document.getElementById('productsContainer');
    iframe.width = "1300px";
    // iframe.width = "100%";
    iframe.height = iframe.contentWindow.document.body.scrollHeight;
}

function saveIframeScreenshot(curCategoryKey, curSubcategory, productUrl, callback) {
    var body = $(document.querySelector("#productsContainer")).contents().find('body')[0];
    html2canvas(body, { letterRendering: true }).then(canvas => {
        exportProductImage(curCategoryKey, curSubcategory, productUrl, canvas);
        callback();
    });
}

function exportProductImage(curCategoryKey, curSubcategory, productUrl, canvas) {
    var imgData = canvas.toDataURL('image/png');

    var prodId = productUrl.split("?id=")[1].split("&")[0];

    var d = new Date();
    var month = d.getMonth() + 1;
    var day = d.getDate();
    var todayStr = d.getFullYear() + '_' +
        (month < 10 ? '0' : '') + month + '_' +
        (day < 10 ? '0' : '') + day;

    var fileName = todayStr + "_" + curCategoryKey + "_" + curSubcategory + "_" + prodId + '.png';
    var a = document.createElement("a"); //Create <a>
    a.href = imgData; //Image Base64 Goes here
    a.download = fileName; //File name Here
    a.click(); //Downloaded file
}

function htmlDecode(input) {
    var doc = new DOMParser().parseFromString(input, "text/html");
    return doc.documentElement.textContent;
}

function updateimagesProgressBar(totalProducts, processedProducts) {
    var percentage = (processedProducts / totalProducts) * 100;

    if (percentage > 100) percentage = 100;
    $('#imagesProgressBar').css('width', percentage + '%');
    $('#imagesProgressBar').html(processedProducts + " / " + totalProducts);
}


function checkFestive(){
    var dateFrom = "11/12";
    var dateTo = "31/12";
    var dateCheck = currentDate();

    var d1 = dateFrom.split("/");
    var d2 = dateTo.split("/");
    var c = dateCheck.split("/");

    var from = new Date(2021, parseInt(d1[1]) - 1, d1[0]);  // -1 because months are from 0 to 11
    var to = new Date(2021, parseInt(d2[1]) - 1, d2[0]);
    var check = new Date(2021, parseInt(c[1]) - 1, c[0]);

    if (check > from && check < to) {
        $("#nav-logo-container").html("Catalog Generatoh oh oh 🎄🎁🎅")
    }
}

function currentDate() {
    var d = new Date();
    var month = d.getMonth() + 1;
    var day = d.getDate();

    return  (day < 10 ? '0' : '') + day+ '/' + (month < 10 ? '0' : '') + month ;
}