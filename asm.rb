require 'byebug'

module ASM
  INSTRUCTIONS = {
    add: 0b10000000,
    shr: 0b10010000,
    shl: 0b10100000,
    not: 0b10110000,
    and: 0b11000000,
    or: 0b11010000,
    xor: 0b11100000,
    cmp: 0b11110000,
  
    #数据
    load: 0b00000000,
    store: 0b00010000,
    data: 0b00100000,
  
    # 跳转
    jmpr: 0b00110000,
    jmp: 0b01000000,
  
    jc: 0b01011000,
    ja: 0b01010100,
    je: 0b01010010,
    jz: 0b01010001,
    jca: 0b01011100,
    jce: 0b01011010,
    jcz: 0b01011001,
    jae: 0b01010110,
    jaz: 0b01010101,
    jez: 0b01010011,
    jcae: 0b01011110,
    jcaz: 0b01011101,
    jcez: 0b01011011,
    jaez: 0b01010111,
    jcaez: 0b01011111,
  
    # 清除
    clf: 0b01100000,
  
    # IN
    in: 0b01110000,
    out: 0b01111000
  }

	# 寄存器
	REG0 = 0b00
	REG1 = 0b01
	REG2 = 0b10
	REG3 = 0b11
end

# 定义汇编语言的语法规则
module AssemblyGrammar
  # 指令 todo
  INSTRUCTIONS = ASM::INSTRUCTIONS.keys
  # 寄存器
  REGISTERS = %w(r0 r1 r2 r3)
  # 标签
  LABEL_REGEX = /(\w+):/
  # 常量
  CONSTS = /(%[a-z]+[a-z0-9_]*)\s*=\s*(0x[0-9a-f]+|0b[0-1]+|\d+)/

  # 操作数
  OPERAND_REGEX = /([a-z]+\d?|\d+|0x[0-9a-f]+|0b[0-1]+|%[a-z]+[a-z0-9_]*)/
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
      elsif match = line.match(CONSTS)
        @tokens << [:CONST, [match[1], match[2]]]
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
            if %w[in out].include?(instruction)
              raise "in/out only receive data/addr" unless %w[data addr].include?(operand)
              [:OPERAND,operand == 'data' ? 0b0000 : 0b0100]
            end
            [:OPERAND, operand]
          end
        end
        @tokens << [:INSTRUCTION, instruction, operands]
        if 'data' == instruction
          @tokens << [:DATA, operands[1][1],[] ]
        elsif 'jmp' == instruction
          @tokens << [:LABEL_DATA, operands[0][1],[] ]
        end
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
    @consts = {}
    @instructions = []
  end
  def parse
    # 第一遍扫描，收集标签的位置信息
    # addition
    # addition = 0
    index = 0
    @tokens.each do |token|
      # if token[0] == :DATA
      #   addition += 1
      # els
      if token[0] == :LABEL
        @labels[token[1]] = index
      elsif token[0] == :CONST
        @consts[token[1][0]] = eval(token[1][1])
      else
        index += 1
      end
    end
    # 第二遍扫描，将汇编代码转换成中间代码
    @tokens.each do |token|
      case token[0]
      when :INSTRUCTION
        instruction = token[1]
        operands = token[2].map do |operand|
          if operand[0] == :REGISTER
            register_index(operand[1])
          else
            # label 或者 操作数(10进制,2进制，16进制)
            lab = @labels[operand[1]]

            # 常量
            cst = @consts[operand[1]]
            # label
            if lab
              lab
            elsif cst
              cst
            # 寄存器
            elsif ri = register_index(operand[1])
              ri
            # 操作数
            elsif /^\d+|0x[a-f0-9]+|0b[0-1]+$/ =~ operand[1]
              eval(operand[1])
            # 其他操作数
            else
              eval(operand[1])
            end
          end
        end
        @instructions << [instruction, operands]
      when :DATA
        @instructions << ['line_data', @consts[token[1]] || eval(token[1])]
      when :LABEL_DATA
        @instructions << ['line_data', @labels[token[1]]]
      end
    end
    @instructions
  end
  private
  def register_index(register)
    case register
    when 'r0'
      ASM::REG0
    when 'r1'
      ASM::REG1
    when 'r2'
      ASM::REG2
    when 'r3'
      ASM::REG3
    end
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
      when 'add','shr','shl','not','and','or','xor','cmp','load','store'
        @code << (ASM::INSTRUCTIONS[instruction[0].to_sym] | instruction[1][0] << 2 | instruction[1][1])
      when 'data'
        @code << (ASM::INSTRUCTIONS[:data] | instruction[1][0])
      when 'jc','ja','je','jz','jca','jce','jcz','jez','jae','jaz','jez','jcae','jcaz','jcez','jaez','jcaez','clf'
        @code << ASM::INSTRUCTIONS[instruction[0].to_sym]
      when 'jmpr'
        @code << (ASM::INSTRUCTIONS[:jmpr] | instruction[1][0])
      when 'jmp'
        @code << ASM::INSTRUCTIONS[:jmp]
      when 'line_data'
        @code << instruction[1]
      when 'in', 'out'
        @code << (ASM::INSTRUCTIONS[instruction[0].to_sym] | instruction[1][0] | instruction[1][1])
      end
    end
    @code
  end
end
# 测试代码
code = <<~CODE
  %abc = 0xaf
  start:
    data r0, %abc
    data r1, 2
    add r0, r1
    cmp r0, r1
    jmp end
  end:
    add r0, r1
CODE
lexer = Lexer.new(code)
tokens = lexer.tokenize
parser = Parser.new(tokens)
instructions = parser.parse
generator = CodeGenerator.new(instructions)
code = generator.generate
puts code.map { |c| "%02x" % c }.join(' ')
# 大端代码