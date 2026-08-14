package view

import tele "gopkg.in/telebot.v4"

func MainMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	btnBuy := menu.Text("🛒 Купить подписку")
	btnMy := menu.Text("🛒 Купить подписку")
	btnSupport := menu.Text("🛒 Купить подписку")	
	menu.Reply(
		menu.Row(btnBuy),
		menu.Row(btnMy, btnSupport),
	)

	return menu
}

func PlansKeyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}
	menu.Inline(menu.Row(menu.Data("1 месяц — 199₽", "plan", "1m")),
		menu.Row(menu.Data("3 месяца — 499₽", "plan", "3m")),
		menu.Row(menu.Data("6 месяцев — 899₽", "plan", "6m")),
		menu.Row(menu.Data("1 год — 1490₽", "plan", "12m")),
		menu.Row(menu.Data("« Назад", "back_main")),
	)
	return menu
}

func LocationsKeyboard(selected map[string]bool) *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}
	locs := []struct{ID, name string} {
		{"nl", "🇳🇱 Нидерланды"},
		{"de", "🇩🇪 Германия"},
		{"nw", "nw Норвегия"},
		{"bg", "bg Бельгия"},
	}

	var rows []tele.Row
	for _, loc := range locs {
		text := loc.name
		if selected[loc.ID] {
			text = "✅ " + text
		}
		rows = append(rows, menu.Row(menu.Data(text, "loc", loc.ID)))
	}

	rows = append(rows,
		menu.Row(menu.Data("Продолжить →", "loc_done")),
		menu.Row(menu.Data("« Назад", "back_plans")),
	)
	menu.Inline(rows...)
	return menu
}

func ConfirmKeyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(menu.Data("💳 Перейти к оплате", "pay")),
		menu.Row(menu.Data("« Назад", "back_locs")),
	)

	return menu
}