<?php

require_once __DIR__ . '/vendor/autoload.php';

$faker = Faker\Factory::create('en');

$timestamp = new \Google\Protobuf\Timestamp();
$timestamp->fromDateTime($faker->dateTimeBetween('-7 days'));

$numbOfItems = $faker->numberBetween(1, 4);
$items = [];

$order = new \PB\Order();
$order->setNumber(sprintf("LI/PHP/%s", $faker->randomNumber(6)));
$order->setStatus($faker->randomElement([
    \PB\OrderStatus::ORDER_STATUS_NEW,
    \PB\OrderStatus::ORDER_STATUS_ACCEPTED,
    \PB\OrderStatus::ORDER_STATUS_CANCALED,
    \PB\OrderStatus::ORDER_STATUS_PAID,
    \PB\OrderStatus::ORDER_STATUS_SENT,
]));
for ($i = 0; $i < $numbOfItems; $i++) {
    $items[] = (new \PB\OrderItem())
        ->setSku($faker->ean8())
        ->setQuantity($faker->numberBetween(1, 5))
        ->setUnitPrice($faker->numberBetween(199, 499));
}
$order->setItems($items);
$order->setCreatedAt($timestamp);

$binaryString = $order->serializeToString();
file_put_contents('protobufer-message.bin', $binaryString);
