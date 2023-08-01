/*
 * Copyright (C) 2023 Yinzuo Jiang
 */

#pragma once

#include <iostream>

static inline void helloA(int a) {
  int b = a + 12;
  std::cout << "Hello, " << a << std::endl;
  std::cout << b << std::endl;
}
