extends Camera2D

@export var zoom_factor := 0.2  # 缩放因子,控制缩放速度
@export var min_zoom := 0.1     # 最小缩放值
@export var max_zoom := 5.0     # 最大缩放值
var is_dragging := false # 是否正在拖动
var mouse_start_pos := Vector2() # 开始拖动时鼠标的位置
var camera_start_pos := Vector2() # 开始拖动时摄像机的位置

func _input(event):
	# 如果是鼠标按钮事件
	if event is InputEventMouseButton:
		# 如果是鼠标左键
		if event.button_index == MOUSE_BUTTON_LEFT:
			# 如果是按下
			if event.pressed:
				# 设置为正在拖动
				is_dragging = true
				# 记录开始拖动时鼠标的位置
				mouse_start_pos = event.position
				# 记录开始拖动时摄像机的位置
				camera_start_pos = position
			# 如果是释放
			else:
				# 设置为不再拖动
				is_dragging = false
	# 如果是鼠标移动事件并且正在拖动
	elif event is InputEventMouseMotion and is_dragging:
		# 计算新的摄像机位置
		# 新位置 = 开始拖动时的摄像机位置 - (当前鼠标位置 - 开始拖动时的鼠标位置)
		position = camera_start_pos - (event.position - mouse_start_pos)
	if event.is_action_pressed("zoom_up"):
		# 滚轮向上滚动,缩小缩放范围
		zoom_camera(zoom_factor)
	elif event.is_action_pressed("zoom_down"):
		# 滚轮向下滚动,放大缩放范围
		zoom_camera(-zoom_factor)

func zoom_camera(factor):
	# 获取当前缩放值
	var current_zoom = get_zoom()
	
	# 计算新的缩放值
	var new_zoom = current_zoom + Vector2(factor, factor)
	
	# 确保新的缩放值在允许范围内
	new_zoom.x = clamp(new_zoom.x, min_zoom, max_zoom)
	new_zoom.y = clamp(new_zoom.y, min_zoom, max_zoom)
	
	# 设置新的缩放值
	set_zoom(new_zoom)
	
	
