package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedProducts(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO products (name, category_id, stock, price, description)
		VALUES
			('Espresso',             1, FLOOR(RANDOM() * 6) + 5, 25000, 'Pure coffee with a strong, concentrated flavor and rich crema'),
			('Cappuccino',           1, FLOOR(RANDOM() * 6) + 5, 30000, 'Espresso with steamed milk and thick milk foam, dusted with cocoa powder'),
			('Latte',                1, FLOOR(RANDOM() * 6) + 5, 32000, 'Espresso with more steamed milk and light foam for a creamier texture'),
			('Americano',            1, FLOOR(RANDOM() * 6) + 5, 22000, 'Espresso diluted with hot water for a smoother coffee experience'),
			('Mocha',                1, FLOOR(RANDOM() * 6) + 5, 35000, 'Espresso combined with chocolate syrup and steamed milk, topped with whipped cream'),
			('Matcha Latte',         2, FLOOR(RANDOM() * 6) + 5, 28000, 'Premium green tea powder whisked with steamed milk for a creamy texture'),
			('Chai Tea',             2, FLOOR(RANDOM() * 6) + 5, 27000, 'Spiced tea with cinnamon, cardamom, and ginger mixed with steamed milk'),
			('Hot Chocolate',        2, FLOOR(RANDOM() * 6) + 5, 25000, 'Rich chocolate drink made with real melted chocolate and steamed milk'),
			('Iced Lemon Tea',       2, FLOOR(RANDOM() * 6) + 5, 20000, 'Refreshing black tea with fresh lemon juice and a hint of sweetness'),
			('Taro Milk Tea',        2, FLOOR(RANDOM() * 6) + 5, 30000, 'Creamy purple taro root powder blended with milk for a unique flavor'),
			('Avocado Toast',        3, FLOOR(RANDOM() * 6) + 5, 45000, 'Sourdough bread topped with mashed avocado, cherry tomatoes, and feta cheese'),
			('Egg Sandwich',         3, FLOOR(RANDOM() * 6) + 5, 35000, 'Soft bread with scrambled eggs, cheese, and fresh vegetables'),
			('Chicken Panini',       3, FLOOR(RANDOM() * 6) + 5, 50000, 'Grilled sandwich with chicken breast, mozzarella, and pesto sauce'),
			('Quiche Lorraine',      3, FLOOR(RANDOM() * 6) + 5, 42000, 'Savory French tart with bacon, cheese, and egg custard in flaky pastry'),
			('Greek Salad',          3, FLOOR(RANDOM() * 6) + 5, 40000, 'Fresh salad with cucumbers, tomatoes, olives, feta cheese, and olive oil'),
			('Tiramisu',             4, FLOOR(RANDOM() * 6) + 5, 38000, 'Classic Italian dessert with layers of coffee-soaked ladyfingers and mascarpone cream'),
			('Chocolate Lava Cake',  4, FLOOR(RANDOM() * 6) + 5, 45000, 'Warm chocolate cake with a molten center, served with vanilla ice cream'),
			('Cheesecake',           4, FLOOR(RANDOM() * 6) + 5, 40000, 'Creamy New York-style cheesecake with graham cracker crust'),
			('Crème Brûlée',         4, FLOOR(RANDOM() * 6) + 5, 42000, 'Rich custard topped with a layer of hardened caramelized sugar'),
			('Macarons',             4, FLOOR(RANDOM() * 6) + 5, 35000, 'French meringue-based cookies with ganache filling in various flavors'),
			('Almond Croissant',     5, FLOOR(RANDOM() * 6) + 5, 28000, 'Buttery croissant filled with almond paste and topped with sliced almonds'),
			('Cinnamon Roll',        5, FLOOR(RANDOM() * 6) + 5, 25000, 'Sweet rolled pastry with cinnamon sugar filling and cream cheese icing'),
			('Chocolate Chip Cookie',5, FLOOR(RANDOM() * 6) + 5, 15000, 'Classic cookie loaded with semi-sweet chocolate chips'),
			('Blueberry Muffin',     5, FLOOR(RANDOM() * 6) + 5, 22000, 'Moist muffin bursting with fresh blueberries and topped with sugar crust'),
			('Banana Bread',         5, FLOOR(RANDOM() * 6) + 5, 20000, 'Homestyle loaf made with ripe bananas and walnuts'),
			('Whipped Cream',        6, FLOOR(RANDOM() * 6) + 5,  5000, 'Freshly whipped heavy cream for topping drinks and desserts'),
			('Caramel Drizzle',      6, FLOOR(RANDOM() * 6) + 5,  7000, 'Sweet caramel sauce for decorating beverages and sweets'),
			('Chocolate Shavings',   6, FLOOR(RANDOM() * 6) + 5,  8000, 'Premium dark chocolate curls for garnishing'),
			('Cinnamon Powder',      6, FLOOR(RANDOM() * 6) + 5,  3000, 'Ground cinnamon for dusting on coffee and baked goods'),
			('Vanilla Syrup',        6, FLOOR(RANDOM() * 6) + 5,  6000, 'Sweet vanilla-flavored syrup for customizing drinks')
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
