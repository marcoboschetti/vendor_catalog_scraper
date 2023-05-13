var categories;

function loadItemsFromFilters() {
    var subCatID = $(".subcategory-filter-item.selected").first().attr("subcat-id");
    var page = parseInt($(".product__pagination a.active").html());

    var loader = `<div class="spinner-border text-info" role="status"><span class="visually-hidden"></span></div>`;
    $("#productsContainer").html(loader);
    $("#searchPagesText").html("Cargando...");
    $(".pagination_container").html("");

    $.get("/api/admin/products/subcategory/" + subCatID + "/" + page, function (data) {
        populateSearchedItems(data.products)
        populatePageSearches(page, data.total, data.page_size);
    });
}

function populatePageSearches(page, total, pageSize) {
    var maxResult = (page * pageSize);
    if(maxResult > total){
        maxResult = total;
    }
    $("#searchPagesText").html("Mostrando resultados " + (((page-1) * pageSize) + 1) + "–" + maxResult + " de " + total);

    var page = parseInt(page);  
    var startPaging = page-1;
    var lastPage = Math.floor(total/pageSize)+1;
    
    var isPageAActive ="";
    var isPageBActive ="active";
    if(page <= 1){
        startPaging = 1;
        isPageBActive ="";
        isPageAActive ="active";
    }

    var pageB = '<a href="#" class="'+isPageBActive+'" >'+(startPaging+1)+'</a>';
    if(startPaging+1 >= lastPage){
        lastPage = "";
        pageC = "";
    }
    var pageC = '<a href="#">'+(startPaging+2)+'</a>';
    if(startPaging+2 >= lastPage){
        pageC = "";
    }
    var endPage = '<a href="#">'+lastPage+'</a>';
    if(startPaging == lastPage){
        endPage = "";
    }

    var elipsis = "<span>...</span>";
    if(startPaging+2 >= lastPage){
        elipsis = "";
    }

    var paginationHTML = `
    <a class="`+isPageAActive+`"  href="#">`+startPaging+`</a>`
    +pageB
    +pageC
    +elipsis
    +endPage
    $(".pagination_container").html(paginationHTML);


    $(".pagination_container a").click(changePage);
}

function changePage(){
    $(".pagination_container a.active").removeClass("active");
    $(this).addClass("active");
    loadItemsFromFilters();
}

function populateSearchedItems(products) {

    var html = "";

    var imagesToLoad = [];
    products.forEach(prod => {
        console.log(prod);
        html += `
        <div class="col-lg-3 col-md-6 col-sm-6" id="`+ prod.id + `">
            <div class="product__item">
                    <div class="product__item__pic set-bg" style="background-image: url('/site/img/loading.gif'); background-size: 100%">
                        <ul class="product__hover">
                        <!--<li><a href="#"><img src="/site/template/img/icon/heart.png" alt=""></a></li>-->
                            <li><a href="/site/item.html?it=`+prod.id+`"><img src="/site/template/img/icon/search.png" alt=""></a></li>
                        </ul>
                    </div>
                <div class="product__item__text">
                    <h6>`+ prod.page_title + `</h6>
                    <div class="prod-subtext">`+ prod.info_block + `</div>
                    <a href="#" class="add-cart">+ Agregar al carrito</a>
                </div>
            </div>
        </div>`

        var bgImg = new Image();
        bgImg.onload = function () {
            $("#" + prod.id + " div div.product__item__pic").first().css("background-image", 'url(' + bgImg.src + ')');
        };
        bgImg.prodID = prod.id;
        bgImg.srcToLoad = "/"+prod.image_urls[0];
        imagesToLoad.push(bgImg)
    });

    $("#productsContainer").html(html);

    imagesToLoad.forEach(bgImg => {
        bgImg.src = bgImg.srcToLoad;
    });
}