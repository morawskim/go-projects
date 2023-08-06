<?php

require_once __DIR__ . '/vendor/autoload.php';

$timestamp = new \Google\Protobuf\Timestamp();
$timestamp->fromDateTime(new DateTime());

$order = new \PB\Order();
$order->setNumber("LI/PHP/123456789");
$order->setStatus(\PB\OrderStatus::ORDER_STATUS_NEW);
$order->setItems([
    (new \PB\OrderItem())->setSku("456789")->setQuantity(2)->setUnitPrice(499),
    (new \PB\OrderItem())->setSku("AZ456")->setQuantity(1)->setUnitPrice(999),
]);
$order->setCreatedAt($timestamp);

$binaryString = $order->serializeToString();
file_put_contents('protobufer-message.bin', $binaryString);
