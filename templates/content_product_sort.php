<div style="min-height: 350px;">
<?php
	include("../../php_lib/common_dao.php");
	//food id of list item
	$product_id = (isset($_GET['product_id'])) ? $_GET['product_id'] : -1;

	if($product_id > -1){
		update("UPDATE orders SET status = 1, finish_time=NOW() WHERE product_id = ".$product_id);
	}
	$products =getSortedByProductOrders();
	?>
	<div class="grid-container">
	<?php
	foreach( $products as $product){
		$message = $product['count']." ".$product['name']." đã hoàn thành!";
	?>
	<div class="order" onclick="onFinishProduct(<?php echo "'".$product['user_id']."',".$product['product_id'].",'".$message."'";?>);">
			<div class="order_title"><span class="special"><?php echo $product['zone'];?></div>
			<div class="order_name"><span class="special"><?php echo '['.$product['count'].'] ';?></span><?php echo $product['name'];?></div>
			<div class="order_comment"><?php if($product['comments']!=null) echo $product['comments'];?></div>
	</div>
<?php } ?>
</div
