extends Node2D

func set_size(that,sprite2D,roche_limit) -> void:
	sprite2D.scale=Vector2(that.mass,that.mass)
	var wh=sprite2D.get_rect().size*sprite2D.scale
	set_collision_shape2D(wh,that)
	set_panel(wh,that)
	
	roche_limit=wh.x/2
	pass

func set_collision_shape2D(wh:Vector2,that):
	# 创建一个 CollisionShape2D 节点
	var collision_shape = CollisionShape2D.new()
	# 创建一个圆形形状
	var shape = CircleShape2D.new()
	shape.radius =wh.x/2
	# 将形状分配给 CollisionShape2D 节点
	collision_shape.shape = shape
	that.add_child(collision_shape)
	
func set_panel(wh:Vector2,that):
	var panel = Panel.new()
	panel.size = wh
	panel.position = -(wh / 2)

	var panel_theme = StyleBoxFlat.new()
	panel_theme.bg_color = Color("#99999900")  # 使用Color类来设置颜色
	panel_theme.set_corner_radius_all(50)

	# 创建一个新的Theme实例
	var theme = Theme.new()

	# 为Panel类型设置StyleBoxFlat样式
	theme.set_stylebox("panel", "Panel", panel_theme)
	# 将Theme应用到panel上
	panel.theme = theme
	
	that.add_child(panel)
	
	# 对于mouse_entered信号
	var callable_entered = Callable(that, "_on_panel_mouse_entered")
	panel.connect("mouse_entered", callable_entered)
	
	# 对于mouse_exited信号
	var callable_exited = Callable(that, "_on_panel_mouse_exited")
	panel.connect("mouse_exited", callable_exited)
