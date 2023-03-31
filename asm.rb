require 'byebug'
# 定义汇编语言的语法规则
module AssemblyGrammar
  # 指令
  INSTRUCTIONS = %w(mov add sub jmp cmp)
  # 寄存器
  REGISTERS = %w(ax bx cx dx)
  # 标签
  LABEL_REGEX = /(\w+):/
  # 操作数
  OPERAND_REGEX = /([a-z]+|\d+)/
  # 操作数列表
  OPERANDS_REGEX = /#{OPERAND_REGEX},\s*#{OPERAND_REGEX}/
  # 完整的指令格式
  INSTRUCTION_REGEX = /^(#{INSTRUCTIONS.join('|')})\s+(#{OPERANDS_REGEX}|#{OPERAND_REGEX})$/
  # 寄存器格式
  REGISTER_REGEX = /^#{REGISTERS.join('|')}$/
end
# 词法分析器
class Lexer
  include AssemblyGrammar
  def initialize(code)
    @code = code
    @tokens = []
  end
  def tokenize
    @code.each_line do |line|
      line = line.strip.downcase
      # 匹配标签
      if (match = line.match(LABEL_REGEX))
        @tokens << [:LABEL, match[1]]
      end
      # 匹配指令
      if (match = line.match(INSTRUCTION_REGEX))
        instruction = match[1]
        operands = match[2].split(',').map(&:strip)
        # 处理寄存器
        operands = operands.map do |operand|
          if operand.match(REGISTER_REGEX)
            [:REGISTER, operand]
          else
            [:OPERAND, operand]
          end
        end
        @tokens << [:INSTRUCTION, instruction, operands]
        byebug
      end
    end
    @tokens
  end
end
# 语法分析器
class Parser
  include AssemblyGrammar
  def initialize(tokens)
    @tokens = tokens
    @labels = {}
    @instructions = []
  end
  def parse
    # 第一遍扫描，收集标签的位置信息
    @tokens.each_with_index do |token, index|
      if token[0] == :LABEL
        @labels[token[1]] = index
      end
    end
    # 第二遍扫描，将汇编代码转换成中间代码
    @tokens.each do |token|
      case token[0]
      when :INSTRUCTION
        instruction = token[1]
        operands = token[2].map do |operand|
          if operand[0] == :REGISTER
            operand[1]
          else
            @labels[operand[1]] || operand[1].to_i
          end
        end
        @instructions << [instruction, operands]
      end
    end
    @instructions
  end
end
# 目标代码生成器
class CodeGenerator
  def initialize(instructions)
    @instructions = instructions
    @code = []
  end
  def generate
    @instructions.each do |instruction|
      case instruction[0]
      when 'mov'
        @code << 0b0001_0000_0000_0000 | (register_index(instruction[1][0]) << 8) | register_index(instruction[1][1])
      when 'add'
        @code << 0b0001_1000_0000_0000 | (register_index(instruction[1][0]) << 8) | register_index(instruction[1][1])
      when 'sub'
        @code << 0b0001_1100_0000_0000 | (register_index(instruction[1][0]) << 8) | register_index(instruction[1][1])
      when 'jmp'
        @code << 0b0100_0000_0000_0000 | (instruction[1] << 8)
      when 'cmp'
        @code << 0b0010_0000_0000_0000 | (register_index(instruction[1][0]) << 8) | register_index(instruction[1][1])
      end
    end
    @code
  end
  private
  def register_index(register)
    case register
    when 'ax'
      0b00
    when 'bx'
      0b01
    when 'cx'
      0b10
    when 'dx'
      0b11
    end
  end
end
# 测试代码
code = <<~CODE
  start:
    mov ax, 1
    mov bx, 2
    add ax, bx
    cmp ax, 10
    jmp end
  end:
    sub ax, bx
CODE
lexer = Lexer.new(code)
tokens = lexer.tokenize
parser = Parser.new(tokens)
instructions = parser.parse
generator = CodeGenerator.new(instructions)
code = generator.generate
puts code.map { |c| "%04x" % c }.join(' ')