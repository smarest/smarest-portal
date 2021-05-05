<div style="min-height: 350px;">
<?php
	include("../../php_lib/common_dao.php");
	//food id of list item
	$order_id = (isset($_GET['order_id'])) ? $_GET['order_id'] : -1;

	if($order_id > -1){
		update("UPDATE orders SET status = 1, finish_time=NOW() WHERE id = ".$order_id);
	}
	?>
	<div class="grid-container">
	<?php
	$products =getSortedByTimeOrders();
	foreach( $products as $product){
		$message = $product['count']." ".$product['name']." đã hoàn thành!";
	?>
	<div class="order" onclick="onFinishOrder(<?php echo "'".$product['user_id']."',".$product['id'].",'".$message."'";?>);">
		<div class="order_title"><strong class="special"><?php echo $product['zone'].': ';?></strong><?php echo $product['table_name'];?></div>
		<div class="order_name"><strong class="special"><?php echo '['.$product['count'].'] ';?></strong><?php echo $product['name'];?></div>
		<div class="order_comment"><?php if($product['comments']!=null) echo $product['comments'];?></div>
	</div>
<?php } ?>
</div>
