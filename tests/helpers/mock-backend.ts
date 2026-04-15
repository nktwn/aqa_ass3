import type { Page, Route } from "@playwright/test";
import {
  qaFeaturedProducts,
  qaProductDetails,
  qaProducts,
  qaSuggestions,
  qaUsers,
  type QaUser,
  type QaUserKey,
} from "../data/ecommerce-data";

type CartProduct = {
  id: number;
  name: string;
  image: string;
  price: number;
  quantity: number;
};

type CartSupplier = {
  id: number;
  name: string;
  total_amount: number;
  delivery_fee: number;
  free_delivery_amount: number;
  order_amount: number;
  product_list: CartProduct[];
};

type CartResponse = {
  customer_id: number;
  total: number;
  suppliers: CartSupplier[];
};

type Address = {
  id?: number;
  description: string;
  street: string;
};

type Order = {
  id: number;
  status: string;
  order_date: string;
  delivery_address: string;
  delivery_comment: string;
  payment_status: string;
  supplier: {
    id: number;
    name: string;
  };
  product_list: CartProduct[];
};

type Contract = {
  id: number;
  content: string;
  status: number;
  supplier_signature?: string;
  customer_signature?: string;
};

type MockApiOptions = {
  authenticatedAs?: QaUserKey;
  initialCart?: CartResponse;
  initialAddresses?: Address[];
  initialOrders?: Order[];
  initialContracts?: Contract[];
  unreadCount?: number;
};

type MockApiState = {
  cart: CartResponse;
  addresses: Address[];
  orders: Order[];
  contracts: Contract[];
};

function json(route: Route, status: number, body: unknown) {
  return route.fulfill({
    status,
    contentType: "application/json; charset=utf-8",
    body: JSON.stringify(body),
  });
}

function findUserByAuthHeader(authorization: string | undefined): QaUser | null {
  if (!authorization?.startsWith("Bearer ")) {
    return null;
  }

  const token = authorization.slice("Bearer ".length);
  return Object.values(qaUsers).find((user) => user.accessToken === token) ?? null;
}

function sortProducts(products: typeof qaProducts, sortBy: string, sortOrder: string) {
  const direction = sortOrder === "asc" ? 1 : -1;
  return [...products].sort((left, right) => {
    if (sortBy === "price") {
      return direction * (left.lowest_product_supplier.price - right.lowest_product_supplier.price);
    }
    if (sortBy === "name") {
      return direction * left.name.localeCompare(right.name);
    }
    return direction * (right.id - left.id);
  });
}

function buildCartState(items: Array<{ productId: number; supplierId: number; quantity: number }>): CartResponse {
  const suppliers = new Map<number, CartSupplier>();

  for (const item of items) {
    const product = qaProducts.find((entry) => entry.id === item.productId);
    if (!product) {
      continue;
    }

    const supplier = product.lowest_product_supplier.supplier;
    const current = suppliers.get(item.supplierId) ?? {
      id: supplier.id,
      name: supplier.name,
      total_amount: 0,
      delivery_fee: supplier.delivery_fee,
      free_delivery_amount: supplier.free_delivery_amount,
      order_amount: supplier.order_amount,
      product_list: [],
    };

    current.product_list.push({
      id: product.id,
      name: product.name,
      image: product.image,
      price: product.lowest_product_supplier.price,
      quantity: item.quantity,
    });
    current.total_amount += product.lowest_product_supplier.price * item.quantity;
    suppliers.set(item.supplierId, current);
  }

  const supplierList = [...suppliers.values()];
  return {
    customer_id: 1,
    total: supplierList.reduce((sum, supplier) => sum + supplier.total_amount, 0),
    suppliers: supplierList,
  };
}

export function buildMockCart(items: Array<{ productId: number; supplierId: number; quantity: number }>) {
  return buildCartState(items);
}

export function buildMockOrder(orderId: number, options?: Partial<Order>): Order {
  const product = qaProducts[0];
  return {
    id: orderId,
    status: "Pending",
    order_date: "2026-04-10T10:30:00Z",
    delivery_address: "Almaty, Satpayev 10",
    delivery_comment: "",
    payment_status: "pending",
    supplier: {
      id: product.lowest_product_supplier.supplier.id,
      name: product.lowest_product_supplier.supplier.name,
    },
    product_list: [
      {
        id: product.id,
        name: product.name,
        image: product.image,
        price: product.lowest_product_supplier.price,
        quantity: 2,
      },
    ],
    ...options,
  };
}

export async function mockEcommerceApi(page: Page, options: MockApiOptions = {}) {
  const state: MockApiState = {
    cart: options.initialCart ?? buildCartState([]),
    addresses: options.initialAddresses ?? [],
    orders: options.initialOrders ?? [],
    contracts: options.initialContracts ?? [],
  };
  let nextAddressId = Math.max(0, ...state.addresses.map((address) => address.id ?? 0)) + 1;
  let nextOrderId = Math.max(5000, ...state.orders.map((order) => order.id)) + 1;
  const unreadCount = options.unreadCount ?? 0;

  await page.route("**/api/**", async (route) => {
    const request = route.request();
    const url = new URL(request.url());
    const method = request.method();
    const path = url.pathname;
    const user = findUserByAuthHeader(request.headers().authorization);
    const authUser = user ?? (options.authenticatedAs ? qaUsers[options.authenticatedAs] : null);

    if (path === "/api/auth/login" && method === "POST") {
      const body = request.postDataJSON() as { phone_number?: string; password?: string };
      const matchingUser = Object.values(qaUsers).find(
        (candidate) =>
          candidate.phoneNumber === body.phone_number &&
          candidate.password === body.password,
      );

      if (!matchingUser) {
        return json(route, 401, { error: "invalid credentials" });
      }

      return json(route, 200, {
        access_token: matchingUser.accessToken,
        refresh_token: matchingUser.refreshToken,
      });
    }

    if (path === "/api/product/list" && method === "GET") {
      const search = url.searchParams.get("search")?.toLowerCase().trim() ?? "";
      const minPrice = Number(url.searchParams.get("min_price") ?? "0");
      const maxPrice = Number(url.searchParams.get("max_price") ?? "0");
      const sortBy = url.searchParams.get("sort_by") ?? "newest";
      const sortOrder = url.searchParams.get("sort_order") ?? "desc";
      const offset = Number(url.searchParams.get("offset") ?? "0");
      const limit = Number(url.searchParams.get("limit") ?? "20");

      const filtered = sortProducts(
        qaProducts.filter((product) => {
          const price = product.lowest_product_supplier.price;
          return (
            (!search || product.name.toLowerCase().includes(search)) &&
            (minPrice === 0 || price >= minPrice) &&
            (maxPrice === 0 || price <= maxPrice)
          );
        }),
        sortBy,
        sortOrder,
      );

      return json(route, 200, {
        product_list: filtered.slice(offset, offset + limit),
        total: filtered.length,
      });
    }

    if (path === "/api/product/featured" && method === "GET") {
      return json(route, 200, {
        product_list: qaFeaturedProducts,
        total: qaFeaturedProducts.length,
      });
    }

    if (path === "/api/product/suggest" && method === "GET") {
      const query = url.searchParams.get("q")?.toLowerCase().trim() ?? "";
      const suggestions = qaSuggestions.filter((item) => item.toLowerCase().includes(query));
      return json(route, 200, { suggestions: suggestions.slice(0, 5) });
    }

    const favoriteMatch = path.match(/^\/api\/product\/(\d+)\/favorite$/);
    if (favoriteMatch) {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }
      if (method === "GET") {
        return json(route, 200, { is_favorite: false });
      }
      return json(route, 200, { status: "ok" });
    }

    const reviewsMatch = path.match(/^\/api\/product\/(\d+)\/reviews$/);
    if (reviewsMatch) {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }
      return json(route, 201, { status: "ok" });
    }

    const productMatch = path.match(/^\/api\/product\/(\d+)$/);
    if (productMatch && method === "GET") {
      const productId = Number(productMatch[1]);
      const detail = qaProductDetails[productId];
      if (!detail) {
        return json(route, 404, { error: "not found" });
      }
      return json(route, 200, { product: detail });
    }

    if (path === "/api/notifications/unread-count" && method === "GET") {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }
      return json(route, 200, { count: unreadCount });
    }

    if (path === "/api/user/role" && method === "GET") {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }
      return json(route, 200, { role: authUser.role });
    }

    if (path === "/api/user/profile" && method === "GET") {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }
      return json(route, 200, {
        user: {
          id: authUser.role + 1,
          name: authUser.name,
          phone_number: authUser.phoneNumber,
          role: authUser.role,
        },
      });
    }

    if (path === "/api/cart/" && method === "GET") {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }
      return json(route, 200, state.cart);
    }

    if (path === "/api/cart/add" && method === "POST") {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }

      const body = request.postDataJSON() as {
        product_id: number;
        supplier_id: number;
        quantity: number;
      };
      const items = state.cart.suppliers.flatMap((supplier) =>
        supplier.product_list.map((product) => ({
          productId: product.id,
          supplierId: supplier.id,
          quantity: product.quantity,
        })),
      );
      const existing = items.find(
        (item) => item.productId === body.product_id && item.supplierId === body.supplier_id,
      );

      if (existing) {
        existing.quantity += body.quantity;
      } else {
        items.push({
          productId: body.product_id,
          supplierId: body.supplier_id,
          quantity: body.quantity,
        });
      }

      state.cart = buildCartState(items);
      return json(route, 200, { status: "ok" });
    }

    if (path === "/api/cart/delete" && method === "DELETE") {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }

      const productId = Number(url.searchParams.get("product_id") ?? "0");
      const supplierId = Number(url.searchParams.get("supplier_id") ?? "0");
      const quantity = Number(url.searchParams.get("quantity") ?? "1");
      const items = state.cart.suppliers.flatMap((supplier) =>
        supplier.product_list.map((product) => ({
          productId: product.id,
          supplierId: supplier.id,
          quantity: product.quantity,
        })),
      );
      const target = items.find(
        (item) => item.productId === productId && item.supplierId === supplierId,
      );

      if (!target) {
        return json(route, 404, { error: "cart item not found" });
      }

      target.quantity -= quantity;
      state.cart = buildCartState(items.filter((item) => item.quantity > 0));
      return json(route, 200, { message: "product deleted from cart" });
    }

    if (path === "/api/cart/clear" && method === "DELETE") {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }
      state.cart = buildCartState([]);
      return json(route, 200, { status: "ok" });
    }

    if (path === "/api/user/address" && method === "GET") {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }
      return json(route, 200, { address_list: state.addresses });
    }

    if (path === "/api/user/address" && method === "POST") {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }

      const body = request.postDataJSON() as { address?: Address };
      if (!body.address?.description || !body.address.street) {
        return json(route, 400, { error: "invalid address" });
      }

      state.addresses.push({
        id: nextAddressId++,
        description: body.address.description,
        street: body.address.street,
      });
      return json(route, 200, { status: "ok" });
    }

    if (path === "/api/cart/checkout" && method === "POST") {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }

      const body = request.postDataJSON() as { address_id?: number };
      if (!body.address_id) {
        return json(route, 400, { error: "address is required" });
      }

      const address = state.addresses.find((entry) => entry.id === body.address_id);
      if (!address) {
        return json(route, 400, { error: "address not found" });
      }

      if (state.cart.suppliers.length === 0) {
        return json(route, 400, { error: "invalid cart" });
      }

      const invalidSupplier = state.cart.suppliers.find(
        (supplier) => supplier.total_amount < supplier.order_amount,
      );
      if (invalidSupplier) {
        return json(route, 400, { error: "invalid cart" });
      }

      const createdOrders = state.cart.suppliers.map((supplier) => ({
        id: nextOrderId++,
        status: "Pending",
        order_date: "2026-04-10T10:30:00Z",
        delivery_address: address.street,
        delivery_comment: "",
        payment_status: "pending",
        supplier: {
          id: supplier.id,
          name: supplier.name,
        },
        product_list: supplier.product_list,
      }));
      state.orders = [...createdOrders, ...state.orders];

      return json(route, 200, {
        checkout_url: "/orders",
        order_id: `PAY-${createdOrders[0].id}`,
      });
    }

    if (path === "/api/order" && method === "GET") {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }
      return json(route, 200, { orders: state.orders });
    }

    const orderByIdMatch = path.match(/^\/api\/order\/(\d+)$/);
    if (orderByIdMatch && method === "GET") {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }
      const order = state.orders.find((entry) => entry.id === Number(orderByIdMatch[1]));
      if (!order) {
        return json(route, 404, { error: "not found" });
      }
      return json(route, 200, order);
    }

    if (path === "/api/order/cancel" && method === "POST") {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }
      const body = request.postDataJSON() as { order_id?: number };
      state.orders = state.orders.map((order) =>
        order.id === body.order_id ? { ...order, status: "Cancelled" } : order,
      );
      return json(route, 200, { message: "order cancelled" });
    }

    if (path === "/api/order/status" && method === "POST") {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }
      if (authUser.role !== 1) {
        return json(route, 401, { error: "only suppliers can update orders" });
      }
      const body = request.postDataJSON() as { order_id?: number; new_status_id?: number };
      const statusMap: Record<number, string> = {
        2: "In Progress",
        3: "Completed",
        4: "Cancelled",
      };
      state.orders = state.orders.map((order) =>
        order.id === body.order_id
          ? { ...order, status: statusMap[body.new_status_id ?? 2] ?? order.status }
          : order,
      );
      return json(route, 200, { message: "order status updated" });
    }

    if (path === "/api/contract" && method === "GET") {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }
      return json(route, 200, state.contracts);
    }

    if (path === "/api/contract/sign" && method === "POST") {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }
      return json(route, 200, { message: "signed" });
    }

    if (path === "/api/product/admin" && method === "POST") {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }
      if (authUser.role !== 2) {
        return json(route, 403, { err: "admin access required" });
      }
      return json(route, 201, { id: 9001 });
    }

    const adminProductMatch = path.match(/^\/api\/product\/admin\/(\d+)$/);
    if (adminProductMatch && (method === "PUT" || method === "DELETE")) {
      if (!authUser) {
        return json(route, 401, { error: "unauthorized" });
      }
      if (authUser.role !== 2) {
        return json(route, 403, { err: "admin access required" });
      }
      return json(route, 200, { status: "ok" });
    }

    return route.fulfill({
      status: 501,
      contentType: "application/json; charset=utf-8",
      body: JSON.stringify({ error: `No mock configured for ${method} ${path}` }),
    });
  });

  return state;
}
