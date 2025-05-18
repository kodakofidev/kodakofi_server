package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// func SeedProducts(ctx context.Context, db *pgxpool.Pool) error {
// 	query := `
// 		INSERT INTO products (name, category_id, price, description, is_deleted)
// 		VALUES
// 			('Espresso',             1, 25000, 'Pure coffee with a strong, concentrated flavor and rich crema', false),
// 			('Cappuccino',           1, 30000, 'Espresso with steamed milk and thick milk foam, dusted with cocoa powder', false),
// 			('Latte',                1, 32000, 'Espresso with more steamed milk and light foam for a creamier texture', false),
// 			('Americano',            1, 22000, 'Espresso diluted with hot water for a smoother coffee experience', false),
// 			('Mocha',                1, 35000, 'Espresso combined with chocolate syrup and steamed milk, topped with whipped cream', false),
// 			('Matcha Latte',         2, 28000, 'Premium green tea powder whisked with steamed milk for a creamy texture', false),
// 			('Chai Tea',             2, 27000, 'Spiced tea with cinnamon, cardamom, and ginger mixed with steamed milk', false),
// 			('Hot Chocolate',        2, 25000, 'Rich chocolate drink made with real melted chocolate and steamed milk', false),
// 			('Iced Lemon Tea',       2, 20000, 'Refreshing black tea with fresh lemon juice and a hint of sweetness', false),
// 			('Taro Milk Tea',        2, 30000, 'Creamy purple taro root powder blended with milk for a unique flavor', false),
// 			('Avocado Toast',        3, 45000, 'Sourdough bread topped with mashed avocado, cherry tomatoes, and feta cheese', false),
// 			('Egg Sandwich',         3, 35000, 'Soft bread with scrambled eggs, cheese, and fresh vegetables', false),
// 			('Chicken Panini',       3, 50000, 'Grilled sandwich with chicken breast, mozzarella, and pesto sauce', false),
// 			('Quiche Lorraine',      3, 42000, 'Savory French tart with bacon, cheese, and egg custard in flaky pastry', false),
// 			('Greek Salad',          3, 40000, 'Fresh salad with cucumbers, tomatoes, olives, feta cheese, and olive oil', false),
// 			('Tiramisu',             4, 38000, 'Classic Italian dessert with layers of coffee-soaked ladyfingers and mascarpone cream', false),
// 			('Chocolate Lava Cake',  4, 45000, 'Warm chocolate cake with a molten center, served with vanilla ice cream', false),
// 			('Cheesecake',           4, 40000, 'Creamy New York-style cheesecake with graham cracker crust', false),
// 			('Crème Brûlée',         4, 42000, 'Rich custard topped with a layer of hardened caramelized sugar', false),
// 			('Macarons',             4, 35000, 'French meringue-based cookies with ganache filling in various flavors', false),
// 			('Almond Croissant',     5, 28000, 'Buttery croissant filled with almond paste and topped with sliced almonds', false),
// 			('Cinnamon Roll',        5, 25000, 'Sweet rolled pastry with cinnamon sugar filling and cream cheese icing', false),
// 			('Chocolate Chip Cookie',5, 15000, 'Classic cookie loaded with semi-sweet chocolate chips', false),
// 			('Blueberry Muffin',     5, 22000, 'Moist muffin bursting with fresh blueberries and topped with sugar crust', false),
// 			('Banana Bread',         5, 20000, 'Homestyle loaf made with ripe bananas and walnuts', false),
// 			('Whipped Cream',        6,  5000, 'Freshly whipped heavy cream for topping drinks and desserts', false),
// 			('Caramel Drizzle',      6,  7000, 'Sweet caramel sauce for decorating beverages and sweets', false),
// 			('Chocolate Shavings',   6,  8000, 'Premium dark chocolate curls for garnishing', false),
// 			('Cinnamon Powder',      6,  3000, 'Ground cinnamon for dusting on coffee and baked goods', false),
// 			('Vanilla Syrup',        6,  6000, 'Sweet vanilla-flavored syrup for customizing drinks', false),
// 			('Ice Cube',             6,  3000, 'ice cubes to make your drink cooler.', false),
// 			('Extra Ice Cube',       6,  3000, 'Extra ice cubes to make your drink cooler.', true)
// 		ON CONFLICT (name) DO NOTHING;
// 	`

// 	_, err := db.Exec(ctx, query)
// 	if err != nil {
// 		log.Printf("Failed to seed products: %v", err)
// 		return err
// 	}

// 	log.Println("Seeded products successfully.")
// 	return nil
// }

func SeedProducts(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO public.products (id,"name",category_id,description,price,is_deleted,created_at,updated_at)
		VALUES
			('4e841656-596c-434d-b5bb-1f27f5d7418c'::uuid,'Espresso',1,'Pure coffee with a strong, concentrated flavor and rich crema',25000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('5740f5fe-5178-4933-9e7b-63cab72fd79a'::uuid,'Cappuccino',1,'Espresso with steamed milk and thick milk foam, dusted with cocoa powder',30000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('e6622f29-2a8b-4110-9a82-8616eed29570'::uuid,'Latte',1,'Espresso with more steamed milk and light foam for a creamier texture',32000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('9a7b950f-3664-4df3-9da5-249e91de2b31'::uuid,'Americano',1,'Espresso diluted with hot water for a smoother coffee experience',22000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('8afc2e72-ef45-45b0-a936-62d34bd626bf'::uuid,'Mocha',1,'Espresso combined with chocolate syrup and steamed milk, topped with whipped cream',35000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('7af264d9-fa31-45a4-8948-b2db4c267fd6'::uuid,'Matcha Latte',2,'Premium green tea powder whisked with steamed milk for a creamy texture',28000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('f866f4f6-f89c-4395-90ca-241dfb52951c'::uuid,'Chai Tea',2,'Spiced tea with cinnamon, cardamom, and ginger mixed with steamed milk',27000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('287ec09d-928c-4562-9f29-86ad95dce6f6'::uuid,'Drink Chocolate',2,'Rich chocolate drink made with real melted chocolate and steamed milk',25000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('95a70a1a-10c8-4a3f-80ad-6c430b74ef3e'::uuid,'Lemon Tea',2,'Refreshing black tea with fresh lemon juice and a hint of sweetness',20000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('40425c10-f932-4b44-97a6-681b56a5ddfa'::uuid,'Taro Milk Tea',2,'Creamy purple taro root powder blended with milk for a unique flavor',30000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('0076aee4-9db2-4941-a69d-e07ff562dc3b'::uuid,'Avocado Toast',3,'Sourdough bread topped with mashed avocado, cherry tomatoes, and feta cheese',45000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('aead3bb8-eaa3-4408-a192-cc36b227f464'::uuid,'Egg Sandwich',3,'Soft bread with scrambled eggs, cheese, and fresh vegetables',35000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('a2b74af4-06cc-4004-989e-6150af06926c'::uuid,'Chicken Panini',3,'Grilled sandwich with chicken breast, mozzarella, and pesto sauce',50000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('dfd726cf-8c5c-44ac-8d06-82378bb4c31c'::uuid,'Quiche Lorraine',3,'Savory French tart with bacon, cheese, and egg custard in flaky pastry',42000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('31b9935a-7bd3-4ae0-898d-386e4cffb82e'::uuid,'Greek Salad',3,'Fresh salad with cucumbers, tomatoes, olives, feta cheese, and olive oil',40000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('06bd407e-5fab-4590-9e3f-7dc442af3b42'::uuid,'Tiramisu',4,'Classic Italian dessert with layers of coffee-soaked ladyfingers and mascarpone cream',38000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('978feadd-1c68-479f-a99a-b831b732464b'::uuid,'Chocolate Lava Cake',4,'Warm chocolate cake with a molten center, served with vanilla ice cream',45000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('fad473ac-7dbe-470a-af07-836726e9b1a6'::uuid,'Cheesecake',4,'Creamy New York-style cheesecake with graham cracker crust',40000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('e07cf18a-b204-455d-95e3-eaf489a805b2'::uuid,'Crème Brûlée',4,'Rich custard topped with a layer of hardened caramelized sugar',42000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('9c4817f0-455e-415f-9c79-8a7e6f4fc1ab'::uuid,'Macarons',4,'French meringue-based cookies with ganache filling in various flavors',35000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('40cb38bd-ba7f-4d40-b5e5-0a013b45e4c8'::uuid,'Almond Croissant',5,'Buttery croissant filled with almond paste and topped with sliced almonds',28000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('97db8997-77ad-4793-8a83-e030ff84d4dd'::uuid,'Cinnamon Roll',5,'Sweet rolled pastry with cinnamon sugar filling and cream cheese icing',25000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('d6a142e4-9960-444f-9cca-041ceb595a9a'::uuid,'Chocolate Chip Cookie',5,'Classic cookie loaded with semi-sweet chocolate chips',15000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('9577276b-d4cb-4de2-a371-c4a25e790c63'::uuid,'Blueberry Muffin',5,'Moist muffin bursting with fresh blueberries and topped with sugar crust',22000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('042875a5-da48-456a-a6e7-b6746f43ab02'::uuid,'Banana Bread',5,'Homestyle loaf made with ripe bananas and walnuts',20000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('fb374b9a-ddde-4a0a-8739-325dcf9543dc'::uuid,'Whipped Cream',6,'Freshly whipped heavy cream for topping drinks and desserts',5000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('ed6d3a40-7760-45cd-a8df-106daeef0227'::uuid,'Caramel Drizzle',6,'Sweet caramel sauce for decorating beverages and sweets',7000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('c504c3ea-af18-4bd8-a2e9-d7e773b1ea5d'::uuid,'Chocolate Shavings',6,'Premium dark chocolate curls for garnishing',8000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('986f09ea-9843-45a3-91a2-da3193c12a63'::uuid,'Cinnamon Powder',6,'Ground cinnamon for dusting on coffee and baked goods',3000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('4c8c4d2e-6f9a-4c91-9f5d-0660ce5c3e6e'::uuid,'Vanilla Syrup',6,'Sweet vanilla-flavored syrup for customizing drinks',6000,false,'2025-05-18 18:24:47.642723+07',NULL),
			('8511cb06-1612-4b5a-9ef5-143abdc2077a'::uuid,'Ice Cube',6,'ice cubes to make your drink cooler.',3000,true,'2025-05-18 18:24:47.642723+07',NULL),
			('40a2a975-9223-4a0e-9379-cb055fe8ae98'::uuid,'Extra Ice Cube',6,'Extra ice cubes to make your drink cooler.',3000,false,'2025-05-18 18:24:47.642723+07',NULL)
		ON CONFLICT (name) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed products: %v", err)
		return err
	}

	log.Println("Seeded products successfully.")
	return nil
}
