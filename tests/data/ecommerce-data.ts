export type QaUserKey = "customer" | "supplier" | "admin";

export type QaUser = {
  key: QaUserKey;
  role: 0 | 1 | 2;
  name: string;
  phoneNumber: string;
  password: string;
  accessToken: string;
  refreshToken: string;
};

export type QaProduct = {
  id: number;
  name: string;
  image: string;
  rating_average: number;
  rating_count: number;
  lowest_product_supplier: {
    price: number;
    sell_amount: number;
    supplier: {
      id: number;
      name: string;
      order_amount: number;
      free_delivery_amount: number;
      delivery_fee: number;
    };
  };
};

export const qaUsers: Record<QaUserKey, QaUser> = {
  customer: {
    key: "customer",
    role: 0,
    name: "QA Customer",
    phoneNumber: "+77000000001",
    password: "password1@",
    accessToken: "mock-access-token-customer",
    refreshToken: "mock-refresh-token-customer",
  },
  supplier: {
    key: "supplier",
    role: 1,
    name: "QA Supplier",
    phoneNumber: "+1234567101",
    password: "password1@",
    accessToken: "mock-access-token-supplier",
    refreshToken: "mock-refresh-token-supplier",
  },
  admin: {
    key: "admin",
    role: 2,
    name: "QA Admin",
    phoneNumber: "+79999999999",
    password: "password1@",
    accessToken: "mock-access-token-admin",
    refreshToken: "mock-refresh-token-admin",
  },
};

export const qaProducts: QaProduct[] = [
  {
    id: 101,
    name: "Arabica Coffee Beans 1kg",
    image: "https://example.com/images/coffee.png",
    rating_average: 4.8,
    rating_count: 18,
    lowest_product_supplier: {
      price: 18000,
      sell_amount: 1,
      supplier: {
        id: 201,
        name: "Roastery One",
        order_amount: 30000,
        free_delivery_amount: 40000,
        delivery_fee: 2500,
      },
    },
  },
  {
    id: 102,
    name: "Premium Green Tea Set",
    image: "https://example.com/images/tea.png",
    rating_average: 4.5,
    rating_count: 11,
    lowest_product_supplier: {
      price: 12000,
      sell_amount: 1,
      supplier: {
        id: 202,
        name: "Tea Republic",
        order_amount: 20000,
        free_delivery_amount: 30000,
        delivery_fee: 1800,
      },
    },
  },
  {
    id: 103,
    name: "Ceramic Mug Studio Edition",
    image: "https://example.com/images/mug.png",
    rating_average: 4.2,
    rating_count: 6,
    lowest_product_supplier: {
      price: 7000,
      sell_amount: 1,
      supplier: {
        id: 201,
        name: "Roastery One",
        order_amount: 30000,
        free_delivery_amount: 40000,
        delivery_fee: 2500,
      },
    },
  },
  {
    id: 104,
    name: "Cold Brew Bottle Kit",
    image: "https://example.com/images/bottle.png",
    rating_average: 4.9,
    rating_count: 22,
    lowest_product_supplier: {
      price: 16000,
      sell_amount: 1,
      supplier: {
        id: 203,
        name: "Barista Lab",
        order_amount: 25000,
        free_delivery_amount: 35000,
        delivery_fee: 2000,
      },
    },
  },
];

export const qaFeaturedProducts = qaProducts.slice(0, 3);

export const qaProductDetails = qaProducts.reduce<Record<number, unknown>>((acc, product) => {
  acc[product.id] = {
    product: {
      ...product,
    },
    suppliers: [
      {
        price: product.lowest_product_supplier.price,
        sell_amount: product.lowest_product_supplier.sell_amount,
        supplier: product.lowest_product_supplier.supplier,
      },
    ],
    reviews: [
      {
        id: product.id * 10,
        user_id: 1,
        user_name: "QA Reviewer",
        rating: Math.round(product.rating_average),
        comment: `Review for ${product.name}`,
        created_at: "2026-04-10T09:00:00Z",
      },
    ],
  };
  return acc;
}, {});

export const qaSuggestions = qaProducts.map((product) => product.name);
